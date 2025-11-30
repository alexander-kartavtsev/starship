package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/alexander-kartavtsev/starship/order/internal/config"
	"github.com/alexander-kartavtsev/starship/platform/pkg/closer"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
	httpMiddleware "github.com/alexander-kartavtsev/starship/platform/pkg/middleware/http"
	customMiddleware "github.com/alexander-kartavtsev/starship/shared/pkg/middleware"
)

type App struct {
	diContainer diContainer
	httpServer  *http.Server
	router      *chi.Mux
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	// Канал для ошибок от компонентов
	errChan := make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// HTTP сервер
	go func() {
		if err := a.runHttpServer(ctx); err != nil {
			errChan <- errors.Errorf("grpc server crashed: %v", err)
		}
	}()

	// Консьюмер
	go func() {
		if err := a.runOrderConsumer(ctx); err != nil {
			errChan <- errors.Errorf("OrderConsumer crashed: %v", err)
		}
	}()

	// Ожидание либо ошибки, либо завершения контекста (например, сигнал SIGINT/SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errChan:
		logger.Info(ctx, "Component crashed, shutting down", zap.Error(err))
		// Триггерим cancel, чтобы остановить второй компонент
		cancel()
		// Дожидаемся завершения всех задач (если есть graceful shutdown внутри)
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDi,
		a.initLogger,
		a.initCloser,
		a.initRouter,
		a.initServer,
		a.runMigrations,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initDi(_ context.Context) error {
	a.diContainer = *NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initRouter(ctx context.Context) error {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(httpMiddleware.NewAuthMiddleware(a.diContainer.IamClient(ctx)).Handle)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(customMiddleware.RequestLogger)
	router.Mount("/", a.diContainer.OrderServer(ctx))

	a.router = router

	return nil
}

func (a *App) initServer(ctx context.Context) error {
	a.httpServer = &http.Server{
		Addr:              config.AppConfig().Server.Address(),
		Handler:           a.router,
		ReadHeaderTimeout: config.AppConfig().Server.TimeoutRead(),
	}

	closer.AddNamed("Http Server", func(ctx context.Context) error {
		err := a.httpServer.Shutdown(ctx)
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}

func (a *App) runMigrations(ctx context.Context) error {
	err := a.diContainer.Migrator(ctx).Up()
	// err := a.diContainer.migrator.Down()
	// err = a.diContainer.migrator.Down()
	if err != nil {
		logger.Error(ctx, "Ошибка миграции базы данных", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) runHttpServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP-сервер запущен на порту %s", config.AppConfig().Server.Port()))

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) runOrderConsumer(ctx context.Context) error {
	logger.Info(ctx, "Order Kafka consumer запущен")

	err := a.diContainer.OrderConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
