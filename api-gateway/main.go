package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "github.com/yinxi0607/YixiGroceryAPI/api-gateway/docs"
	"github.com/yinxi0607/YixiGroceryAPI/api-gateway/middleware"
	userProto "github.com/yinxi0607/YixiGroceryAPI/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title YixiGroceryAPI
// @version 1.0
// @description API for Yixi Grocery microservices
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Create Gin router
	r := gin.Default()
	r.Use(middleware.Auth())

	// Create gRPC-Gateway mux
	gwMux := runtime.NewServeMux()

	// gRPC connection to user-service
	conn, err := grpc.NewClient("yinxi-user-service:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to user-service: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {

		}
	}(conn)

	// Register UserService handler
	if err := userProto.RegisterUserServiceHandler(context.Background(), gwMux, conn); err != nil {
		log.Fatalf("Failed to register user service handler: %v", err)
	}

	// Mount gRPC-Gateway to Gin
	r.Any("/api/*any", gin.WrapH(gwMux))

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	log.Println("Starting API Gateway on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
