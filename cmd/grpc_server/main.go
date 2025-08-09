package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"runtime"
	"time"

	"github.com/joho/godotenv"

	dao "github.com/en7ka/auth/deploy/postgres/cmd"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/en7ka/auth/pkg/user_v1"
)

const grpcPort = 50051

func main() {
	// Получаем абсолютный путь к текущему файлу
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	// Строим абсолютный путь к файлу .env
	envPath := filepath.Join(basepath, "../../deploy/.env")

	// Загружаем .env
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatal("failed to listen: 50051 ")
	}

	storage, err := dao.InitStorage()
	if err != nil {
		log.Fatal("failed to init storage")
	}
	defer storage.CloseCon()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserAPIServer(s, &server{storage: storage})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	desc.UnimplementedUserAPIServer
	storage dao.PostgresInterface
}

func toTimestampProto(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func (s *server) Get(_ context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	// В вашем proto файле GetRequest имеет только поле id.
	params := dao.GetUserPar{ID: &req.Id}

	userProfile, err := s.storage.GetUser(params)
	if err != nil {
		return nil, fmt.Errorf("error when getting the user profile: %w", err)
	}

	// Заполнение полей из userProfile в GetResponse
	response := &desc.GetResponse{
		Id:        userProfile.ID,
		Name:      userProfile.Username,
		Email:     userProfile.Email,
		Role:      userProfile.Role,
		CreatedAt: toTimestampProto(userProfile.CreatedAt),
		UpdatedAt: toTimestampProto(userProfile.UpdatedAt),
	}

	return response, nil
}

func (s *server) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	// В вашем proto файле CreateRequest не имеет вложенной User
	user := dao.User{
		Username: req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     desc.Role_user, // Предполагаем, что при создании роль всегда 'user'
	}

	id, err := s.storage.Save(user)
	if err != nil {
		return nil, fmt.Errorf("error when saving the user: %w", err)
	}
	return &desc.CreateResponse{Id: id}, nil
}

func (s *server) Update(_ context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	// В вашем proto файле поля Name и Email - это StringValue
	updateUser := dao.UpdateUser{
		ID:       req.GetId(),
		Username: req.GetName().GetValue(),
		Email:    req.GetEmail().GetValue(),
	}
	err := s.storage.Update(updateUser)
	if err != nil {
		return &emptypb.Empty{}, fmt.Errorf("error updating user: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *server) Delete(_ context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	idDel := dao.DeleteID(req.GetId())
	err := s.storage.Delete(idDel)
	if err != nil {
		return &emptypb.Empty{}, fmt.Errorf("error deleting user: %w", err)
	}

	return &emptypb.Empty{}, nil
}
