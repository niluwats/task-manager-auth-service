package main

import (
	"log"
	"net"
	"os"

	"github.com/niluwats/task-manager-auth-service/api/pb"
	"github.com/niluwats/task-manager-auth-service/internal/db"
	"github.com/niluwats/task-manager-auth-service/internal/handlers"
	"github.com/niluwats/task-manager-auth-service/internal/repositories"
	"github.com/niluwats/task-manager-auth-service/internal/service"
	"google.golang.org/grpc"
)

func main() {
	dbClient := db.ConnectDB()

	userRepo := repositories.NewUserRepositoryDB(dbClient)
	service := service.NewUserService(userRepo)
	handler := handlers.NewAuthHandler(service)

	lis, err := net.Listen("tcp", os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatal("Failed to listen on gRPC port : ", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, handler)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("Failed to serve : ", err)
	}
}
