package main

import (
	"log"
	"net"

	userProto "github.com/yinxi0607/YixiGroceryAPI/proto/user"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/config"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/handler"
	"google.golang.org/grpc"
)

func main() {
	// Initialize database
	config.InitDB()

	// Create gRPC server
	srv := grpc.NewServer()

	// Register UserService
	userProto.RegisterUserServiceServer(srv, &handler.UserHandler{})

	// Start gRPC server
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Starting gRPC server on :8081")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
