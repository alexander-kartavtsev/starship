package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/alexander-kartavtsev/starship/payment/internal/config"
	"github.com/alexander-kartavtsev/starship/platform/pkg/closer"
	"github.com/alexander-kartavtsev/starship/platform/pkg/grpc/health"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
	"github.com/alexander-kartavtsev/starship/platform/pkg/tracing"
	paymentV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/payment/v1"
)

type App struct {
	diContainer diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
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
	return a.runGrpcServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDi,
		a.initLogger,
		a.initTracing,
		a.initCloser,
		a.initListener,
		a.initGrpcServer,
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

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(ctx context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().GRPC.Address())
	if err != nil {
		logger.Error(ctx, "failed to listen", zap.Error(err))
		return err
	}

	closer.AddNamed("TCP Listener", func(context.Context) error {
		err = listener.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}
		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initGrpcServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(tracing.UnaryServerInterceptor(config.AppConfig().Tracing.ServiceName())),
	)
	closer.AddNamed("gRPC server", func(context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)
	health.RegisterService(a.grpcServer)
	paymentV1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.PaymentApi(ctx))

	return nil
}

func (a *App) runGrpcServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC server listening on port %s", config.AppConfig().GRPC.Port()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initTracing(ctx context.Context) error {
	err := tracing.InitTracer(ctx, config.AppConfig().Tracing)
	if err != nil {
		return err
	}

	closer.AddNamed("tracer", tracing.ShutdownTracer)

	return nil
}
