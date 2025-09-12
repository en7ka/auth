package app

import (
	"context"
	"log"

	userApi "github.com/en7ka/auth/internal/api/auth"
	authApi "github.com/en7ka/auth/internal/api/user"
	clRedis "github.com/en7ka/auth/internal/client/cache/redis"
	"github.com/en7ka/auth/internal/client/db"
	"github.com/en7ka/auth/internal/client/db/pg"
	"github.com/en7ka/auth/internal/client/db/transaction"
	"github.com/en7ka/auth/internal/closer"
	"github.com/en7ka/auth/internal/config"
	repoAuth "github.com/en7ka/auth/internal/repository/auth"
	repoRedis "github.com/en7ka/auth/internal/repository/redis"
	repoinf "github.com/en7ka/auth/internal/repository/repositoryinterface"
	repoUser "github.com/en7ka/auth/internal/repository/user"
	servAuth "github.com/en7ka/auth/internal/service/auth"
	"github.com/en7ka/auth/internal/service/servinterface"
	servUser "github.com/en7ka/auth/internal/service/user"
	redigo "github.com/gomodule/redigo/redis"
)

type serviceProvider struct {
	pgConfig     config.PGConfig
	grpcConfig   config.GRPCConfig
	redisConfig  config.RedisConfig
	jwtConfig    config.JWTConfig
	accessConfig config.AccessConfig

	dbClient  db.Client
	redisPool *redigo.Pool
	txManager db.TxManager

	userRepository repoinf.UserRepository
	userCache      repoinf.UserCache
	authRepository repoinf.AuthRepository

	userService servinterface.UserService
	authService servinterface.AuthService

	userImpl *userApi.Controller
	authImpl *authApi.Controller
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

func (s *serviceProvider) GetJWTConfig(_ context.Context) config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			log.Fatalf("failed to load jwt config: %v", err)
		}
		s.jwtConfig = cfg
	}

	return s.jwtConfig
}

func (s *serviceProvider) GetAccessConfig(_ context.Context) config.AccessConfig {
	if s.accessConfig == nil {
		cfg, err := config.NewAccessConfig()
		if err != nil {
			log.Fatalf("failed to load access config: %v", err)
		}
		s.accessConfig = cfg
	}

	return s.accessConfig
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

func (s *serviceProvider) GetAuthRepository(ctx context.Context) repoinf.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = repoUser.NewRepository(s.GetDBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) GetUserCache(_ context.Context) repoinf.UserCache {
	if s.userCache == nil {
		redisClient := clRedis.NewClient(s.GetRedisPool(), s.GetRedisConfig())
		s.userCache = repoRedis.NewRedisCache(redisClient)
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

func (s *serviceProvider) GetAuthService(ctx context.Context) servinterface.AuthService {
	if s.authService == nil {
		s.authService = servUser.NewService(
			s.GetAuthRepository(ctx),
			s.GetUserCache(ctx),
			s.GetTxManager(ctx),
			s.GetJWTConfig(ctx),
			s.GetAccessConfig(ctx),
		)
	}

	return s.authService
}
func (s *serviceProvider) GetUserApiController(ctx context.Context) *userApi.Controller {
	if s.userImpl == nil {
		s.userImpl = userApi.NewImplementation(s.GetUserService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) GetAuthApiController(ctx context.Context) *authApi.Controller {
	if s.authImpl == nil {
		s.authImpl = authApi.NewImplementation(s.GetAuthService(ctx))
	}

	return s.authImpl
}
