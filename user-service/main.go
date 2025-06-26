package main

import (
	"github.com/yinxi0607/YixiGroceryAPI/logger"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/handler"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/model"
	pb "github.com/yinxi0607/YixiGroceryAPI/user-service/proto"
	"go-micro.dev/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	// 初始化日志
	if err := logger.Init("logs/user-service.log", "DEBUG"); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	logger.Logger.Info("Starting user-service...")

	// 初始化数据库
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=yourpassword dbname=user_service port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Logger.Fatal("Failed to connect to database: ", err)
	}

	// 自动迁移用户表
	if err := db.AutoMigrate(&model.User{}); err != nil {
		logger.Logger.Fatal("Failed to migrate database: ", err)
	}
	logger.Logger.Info("Database connected and migrated")

	// 创建 Micro 服务
	service := micro.NewService(
		micro.Name("user-service"),
	)

	// 初始化服务
	service.Init()

	// 注册 Handler，传递数据库实例
	err = pb.RegisterUserServiceHandler(service.Server(), handler.New(db))
	if err != nil {
		logger.Logger.Fatal("Failed to register handler: ", err)
	}

	// 运行服务
	if err := service.Run(); err != nil {
		logger.Logger.Fatal("Failed to run service: ", err)
	}
}
