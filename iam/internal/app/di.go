package app

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	authV1API "github.com/vipshark78/microservices-course-homeworks/iam/internal/api/auth/v1"
	userV1API "github.com/vipshark78/microservices-course-homeworks/iam/internal/api/user/v1"
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/config"
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/repository"
	sessionRepository "github.com/vipshark78/microservices-course-homeworks/iam/internal/repository/session"
	userRepository "github.com/vipshark78/microservices-course-homeworks/iam/internal/repository/user"
	"github.com/vipshark78/microservices-course-homeworks/iam/internal/service"
	authService "github.com/vipshark78/microservices-course-homeworks/iam/internal/service/auth"
	userService "github.com/vipshark78/microservices-course-homeworks/iam/internal/service/user"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/cache"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/cache/redis"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/closer"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/migrator"
	userMigrator "github.com/vipshark78/microservices-course-homeworks/platform/pkg/migrator/pg"
	authV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/auth/v1"
	userV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/user/v1"
)

type diContainer struct {
	authV1API         authV1.AuthServiceServer
	userV1API         userV1.UserServiceServer
	authService       service.AuthService
	userService       service.UserService
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	userMigrator      migrator.Migrator
	pgxPool           *pgxpool.Pool

	redisPool   *redigo.Pool
	redisClient cache.RedisClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) UserV1API(ctx context.Context) userV1.UserServiceServer {
	if d.userV1API == nil {
		d.userV1API = userV1API.NewAPI(d.UserService(ctx))
	}

	return d.userV1API
}

func (d *diContainer) AuthV1API(ctx context.Context) authV1.AuthServiceServer {
	if d.authV1API == nil {
		d.authV1API = authV1API.NewAPI(d.AuthService(ctx))
	}

	return d.authV1API
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewUserService(
			d.UserRepository(ctx),
		)
	}
	return d.userService
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewAuthService(
			d.UserRepository(ctx),
			d.SessionRepository(ctx),
		)
	}
	return d.authService
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewRepository(d.PgxPool(ctx))
		d.UserMigrator(d.PgxPool(ctx))
		err := d.userMigrator.Up()
		if err != nil {
			logger.Fatal(ctx, "Ошибка при миграции", zap.Error(err))
		}
	}
	return d.userRepository
}

func (d *diContainer) UserMigrator(pgxPool *pgxpool.Pool) migrator.Migrator {
	if d.userMigrator == nil {
		d.userMigrator = userMigrator.NewMigrator(stdlib.OpenDB(*pgxPool.Config().ConnConfig), config.AppConfig().Postgres.MigrationsDir())
	}
	return d.userMigrator
}

func (d *diContainer) SessionRepository(ctx context.Context) repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionRepository.NewRepository(d.RedisClient(ctx))
	}
	return d.sessionRepository
}

func (d *diContainer) PgxPool(ctx context.Context) *pgxpool.Pool {
	if d.pgxPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			logger.Fatal(ctx, "failed to connect to database", zap.Error(err))
		}

		closer.AddNamed("PgxPool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.pgxPool = pool
	}
	return d.pgxPool
}

func (d *diContainer) RedisClient(ctx context.Context) cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(
			d.RedisPool(ctx),
			logger.Logger(),
			config.AppConfig().Redis.ConnectionTimeout(),
		)
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
	}
	return d.redisPool
}
