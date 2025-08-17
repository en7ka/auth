package app

import (
	"context"
	"github.com/en7ka/auth/internal/closer"
	"github.com/en7ka/auth/internal/config"
	desc "github.com/en7ka/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type App struct {
	serviceProvaider *serviceProvaider
	grpcServer       *grpc.Server
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServerProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	if err := config.Load(".env"); err != nil {
		return err
	}

	return nil
}

func (a *App) initServerProvider(_ context.Context) error {
	a.serviceProvaider = newServiceProvider()

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	desc.RegisterUserAPIServer(a.grpcServer, a.serviceProvaider.GetUserImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on #{a.serviceProvaider.GetGRPCConfig().Address()}")

	list, err := net.Listen("tcp", a.serviceProvaider.GetGRPCConfig().Address())
	if err != nil {
		return err
	}

	if err = a.grpcServer.Serve(list); err != nil {
		return err
	}

	return nil
}
