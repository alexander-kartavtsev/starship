package app

import (
	"context"
	"log"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authApiV1 "github.com/alexander-kartavtsev/starship/iam/internal/api/auth/v1"
	userApiV1 "github.com/alexander-kartavtsev/starship/iam/internal/api/user/v1"
	"github.com/alexander-kartavtsev/starship/iam/internal/config"
	"github.com/alexander-kartavtsev/starship/iam/internal/repository"
	"github.com/alexander-kartavtsev/starship/iam/internal/repository/session"
	"github.com/alexander-kartavtsev/starship/iam/internal/repository/user"
	"github.com/alexander-kartavtsev/starship/iam/internal/service"
	authService "github.com/alexander-kartavtsev/starship/iam/internal/service/auth"
	userService "github.com/alexander-kartavtsev/starship/iam/internal/service/user"
	"github.com/alexander-kartavtsev/starship/platform/pkg/cache"
	"github.com/alexander-kartavtsev/starship/platform/pkg/cache/redis"
	"github.com/alexander-kartavtsev/starship/platform/pkg/closer"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
	migrator "github.com/alexander-kartavtsev/starship/platform/pkg/migrator/pg"
	"github.com/alexander-kartavtsev/starship/platform/pkg/tracing"
	authV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/auth/v1"
	userV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/user/v1"
)

type diContainer struct {
	authApi authV1.AuthServiceServer
	userApi userV1.UserServiceServer

	authService service.AuthService
	userService service.UserService

	authServiceClient authV1.AuthServiceClient

	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository

	dbConn   *pgx.Conn
	dbPool   *pgxpool.Pool
	migrator *migrator.Migrator

	redisPool   *redigo.Pool
	redisClient cache.RedisClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AuthApi(ctx context.Context) authV1.AuthServiceServer {
	if d.authApi == nil {
		d.authApi = authApiV1.NewApi(d.AuthService(ctx))
		logger.Info(ctx, "Инициализация AuthApi")
	}
	return d.authApi
}

func (d *diContainer) UserApi(ctx context.Context) userV1.UserServiceServer {
	if d.userApi == nil {
		d.userApi = userApiV1.NewApi(d.UserService(ctx))
		logger.Info(ctx, "Инициализация UserApi")
	}
	return d.userApi
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewService(d.SessionRepository(ctx), d.UserRepository(ctx))
		logger.Info(ctx, "Инициализация AuthService")
	}
	return d.authService
}

func (d *diContainer) AuthClient(ctx context.Context) authV1.AuthServiceClient {
	if d.authServiceClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().Grpc.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor("iam-service")),
		)
		if err != nil {
			log.Printf("failed to connect: %v\n", err)
		}

		logger.Info(ctx, "Инициализация AuthService gRPC Client")

		closer.AddNamed("AuthService gRPC client", func(ctx context.Context) error {
			return conn.Close()
		})

		d.authServiceClient = authV1.NewAuthServiceClient(conn)
	}
	return d.authServiceClient
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewService(d.UserRepository(ctx))
		logger.Info(ctx, "Инициализация UserService")
	}
	return d.userService
}

func (d *diContainer) SessionRepository(ctx context.Context) repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = session.NewRepository(d.RedisClient(ctx), d.DbPool(ctx))
		logger.Info(ctx, "Инициализация SessionRepository")
	}
	return d.sessionRepository
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = user.NewRepository(d.DbConn(ctx), d.DbPool(ctx))
		logger.Info(ctx, "Инициализация UserRepository")
	}
	return d.userRepository
}

func (d *diContainer) DbConn(ctx context.Context) *pgx.Conn {
	if d.dbConn == nil {
		conn, err := pgx.Connect(ctx, config.AppConfig().Postgres.Uri())
		if err != nil {
			log.Printf("failed to connect to database: %v\n", err)
			panic(err)
		}

		// Проверяем, что соединение с базой установлено
		err = conn.Ping(ctx)
		if err != nil {
			log.Printf("База данных недоступна: %v\n", err)
			panic(err)
		}

		logger.Info(ctx, "Инициализация DbConn")

		closer.AddNamed("Db connection", func(ctx context.Context) error {
			return conn.Close(ctx)
		})

		d.dbConn = conn
	}
	return d.dbConn
}

func (d *diContainer) DbPool(ctx context.Context) *pgxpool.Pool {
	if d.dbPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.Uri())
		if err != nil {
			log.Printf("Ошибка подключения к б/д: %v\n", err)
		}

		logger.Info(ctx, "Инициализация DbPool")

		closer.AddNamed("Db pool connection", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.dbPool = pool
	}
	return d.dbPool
}

func (d *diContainer) RedisClient(ctx context.Context) cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(d.RedisPool(ctx), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
		logger.Info(ctx, "Инициализация RedisClient")
	}
	return d.redisClient
}

func (d *diContainer) RedisPool(ctx context.Context) *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}
		logger.Info(ctx, "Инициализация RedisPool")
	}
	return d.redisPool
}

func (d *diContainer) Migrator(ctx context.Context) *migrator.Migrator {
	if d.migrator == nil {
		d.migrator = migrator.NewMigrator(
			stdlib.OpenDB(*d.DbConn(ctx).Config().Copy()),
			config.AppConfig().Postgres.MigrationsDir(),
		)
	}
	return d.migrator
}
