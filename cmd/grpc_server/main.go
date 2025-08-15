package main

import (
	"context"
	"github.com/en7ka/auth/internal/config"
	desc "github.com/en7ka/auth/pkg/user_v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	userApi "github.com/en7ka/auth/internal/api/auth"
	userRepo "github.com/en7ka/auth/internal/repository/auth"
	userService "github.com/en7ka/auth/internal/service/auth"
)

func main() {
	ctx := context.Background()

	//считываем переменные окружения
	err := config.Load(".env")
	if err != nil {
		log.Fatalf("ошибка к подключению к .env: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("ошибка к подключению с grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("ошибка к подключению к pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Printf("ошибка в прослушивании: %v", err)
	}

	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Printf("ошибка в подключении: %v", err)
	}
	defer pool.Close()

	noteRepo := userRepo.NewRepository(pool)
	noteService := userService.NewService(noteRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserAPIServer(s, userApi.NewImplementation(noteService))

	log.Printf("сервер слушает на порту: %v", grpcConfig.Address())

	if err = s.Serve(lis); err != nil {
		log.Printf("ошибка в прослушивании: %v", err)
	}

}
