package main

import (
	"go-micro.dev/v5/web"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "github.com/yinxi0607/YixiGroceryAPI/api-gateway/docs"
	"github.com/yinxi0607/YixiGroceryAPI/api-gateway/handler"
	"github.com/yinxi0607/YixiGroceryAPI/api-gateway/middleware"
	"google.golang.org/grpc"
)

// @title YixiGroceryAPI
// @version 1.0
// @description API for Yixi Grocery microservices
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 创建 Go-Micro Web 服务
	service := web.NewService(
		web.Name("go.micro.web.api-gateway.yixi"),
		web.Address(":8080"),
	)

	// 初始化 Gin
	r := gin.Default()
	r.Use(middleware.Auth())

	// gRPC 客户端连接
	conn, err := grpc.Dial("yixi-user-service:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 用户服务路由
	userHandler := handler.NewUserHandler(conn)
	r.POST("/api/auth/register", userHandler.Register)
	r.POST("/api/auth/login", userHandler.Login)
	r.GET("/api/users/me", userHandler.GetUserInfo)
	r.POST("/api/users/addresses", userHandler.AddAddress)
	r.PUT("/api/users/addresses/:id", userHandler.UpdateAddress)
	r.DELETE("/api/users/addresses/:id", userHandler.DeleteAddress)
	r.GET("/api/users/addresses", userHandler.GetAddresses)

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 设置 Gin 路由
	service.Handle("/", r)

	// 启动服务
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
