package main

import (
	"context"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/en7ka/auth/pkg/user_v1"

	"github.com/en7ka/auth/internal/config"
	"github.com/en7ka/auth/internal/repository/auth"
	"github.com/en7ka/auth/internal/repository/auth/converter"
	"github.com/en7ka/auth/internal/repository/auth/model"
	repinf "github.com/en7ka/auth/internal/repository/repositoryinterface"
)

type server struct {
	desc.UnimplementedUserAPIServer
	userRepository repinf.UserRepository
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.userRepository.Create(ctx, &model.UserInfo{
		Username: req.GetInfo().GetUsername(),
		Email:    req.GetInfo().GetEmail(),
		Password: req.GetInfo().GetPassword(),
	})
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{Id: id}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	u, err := s.userRepository.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &desc.GetResponse{
		Note: &desc.Note{
			Id: u.Id,
			Info: &desc.NoteInfo{
				Username: u.Info.Username,
				Email:    u.Info.Email,
				Password: u.Info.Password,
				Role:     converter.RoleFromString(u.Role),
			},
			CreatedAt: timestamppb.New(u.CreatedAt),
			UpdatedAt: timestamppb.New(u.UpdatedAt),
		},
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	id := req.GetId()

	info := &model.UserInfo{}
	if v := req.GetInfo().GetUsername(); v != nil {
		info.Username = v.GetValue()
	}
	if v := req.GetInfo().GetEmail(); v != nil {
		info.Email = v.GetValue()
	}

	if err := s.userRepository.Update(ctx, id, info); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if err := s.userRepository.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func main() {
	ctx := context.Background()

	if err := config.Load(".env"); err != nil {
		log.Fatal(err)
	}

	grpcCfg, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatal(err)
	}

	pgCfg, err := config.NewPGConfig()
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", grpcCfg.Address())
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, pgCfg.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := auth.NewRepository(pool)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserAPIServer(s, &server{userRepository: repo})

	log.Printf("grpc: %s", grpcCfg.Address())

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
