package app

import (
	"context"
	"log"

	"github.com/en7ka/auth/internal/api/auth"
	"github.com/en7ka/auth/internal/client/db"
	"github.com/en7ka/auth/internal/client/db/pg"
	"github.com/en7ka/auth/internal/client/db/transaction"
	"github.com/en7ka/auth/internal/closer"
	"github.com/en7ka/auth/internal/config"
	repoAuth "github.com/en7ka/auth/internal/repository/auth"
	repoRedis "github.com/en7ka/auth/internal/repository/redis"
	repoinf "github.com/en7ka/auth/internal/repository/repositoryinterface"
	servAuth "github.com/en7ka/auth/internal/service/auth"
	"github.com/en7ka/auth/internal/service/servinterface"
	redigo "github.com/gomodule/redigo/redis"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	redisConfig   config.RedisConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig

	dbClient  db.Client
	redisPool *redigo.Pool
	txManager db.TxManager

	userRepository repoinf.UserRepository
	userCache      repoinf.UserCache

	userService servinterface.UserService

	userImpl *auth.Controller
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GetPGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to load pg config: %v", err)
		}
		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *serviceProvider) GetGRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to load gRPC config: %v", err)
		}
		s.grpcConfig = cfg
	}
	return s.grpcConfig
}

func (s *serviceProvider) GetRedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to load redis config: %v", err)
		}
		s.redisConfig = cfg
	}
	return s.redisConfig
}

func (s *serviceProvider) GetDBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.GetPGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		if err = cl.DB().Ping(ctx); err != nil {
			log.Fatalf("failed to ping database: %v", err)
		}
		closer.Add(cl.Close)
		s.dbClient = cl
	}
	return s.dbClient
}

func (s *serviceProvider) GetHTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to load http config: %v", err)
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) GetSwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err)
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) GetRedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.GetRedisConfig().MaxIdle(),
			IdleTimeout: s.GetRedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.GetRedisConfig().Address())
			},
		}
	}

	return s.redisPool
}

func (s *serviceProvider) GetTxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.GetDBClient(ctx).DB())
	}
	return s.txManager
}

func (s *serviceProvider) GetUserRepository(ctx context.Context) repoinf.UserRepository {
	if s.userRepository == nil {
		s.userRepository = repoAuth.NewRepository(s.GetDBClient(ctx))
	}
	return s.userRepository
}

func (s *serviceProvider) GetUserCache(ctx context.Context) repoinf.UserCache {
	if s.userCache == nil {
		s.userCache = repoRedis.NewRedisCache(s.GetRedisPool())
	}
	return s.userCache
}

func (s *serviceProvider) GetUserService(ctx context.Context) servinterface.UserService {
	if s.userService == nil {
		s.userService = servAuth.NewService(
			s.GetUserRepository(ctx),
			s.GetUserCache(ctx),
			s.GetTxManager(ctx),
		)
	}
	return s.userService
}

func (s *serviceProvider) GetUserImpl(ctx context.Context) *auth.Controller {
	if s.userImpl == nil {
		s.userImpl = auth.NewImplementation(s.GetUserService(ctx))
	}
	return s.userImpl
}
