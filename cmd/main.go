package main

import (
	"fmt"
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

	startgRPC(handler)
}

func startgRPC(handler handlers.AuthHandlerImpl) {
	port := os.Getenv("GRPC_PORT")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal("Failed to listen on gRPC port : ", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, handler)

	log.Printf("gRPC server started on port %s", port)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("Failed to serve : ", err)
	}

}
