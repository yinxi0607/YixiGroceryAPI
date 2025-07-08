package main

import (
	"github.com/yinxi0607/YixiGroceryAPI/config"
	"github.com/yinxi0607/YixiGroceryAPI/logger"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/handler"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/model"
	pb "github.com/yinxi0607/YixiGroceryAPI/user-service/proto"
	"go-micro.dev/v5"
	"os"
)

func main() {
	// 初始化日志
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "DEBUG"
	}
	if err := logger.Init("logs/user-service.log", logLevel); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	logger.Logger.Info("Starting user-service...")

	// 加载数据库配置
	dbConfig := config.LoadDBConfig("user-service")

	// 初始化数据库
	db, err := config.InitDB(dbConfig)
	if err != nil {
		logger.Logger.Fatal("Failed to initialize database: ", err)
	}

	// 自动迁移用户表
	if err := db.AutoMigrate(&model.User{}); err != nil {
		logger.Logger.Fatal("Failed to migrate database: ", err)
	}
	logger.Logger.Info("Database migrated")

	// 加载 Redis 配置
	redisConfig := config.LoadRedisConfig("user-service")

	// 初始化 Redis
	redisClient, err := config.InitRedis(redisConfig)
	if err != nil {
		logger.Logger.Fatal("Failed to initialize Redis: ", err)
	}

	// 创建 Micro 服务
	service := micro.NewService(
		micro.Name("user-service"),
	)

	// 初始化服务
	service.Init()

	// 注册 Handler，传递数据库和 Redis 实例
	err = pb.RegisterUserServiceHandler(service.Server(), handler.New(db, redisClient))
	if err != nil {
		logger.Logger.Fatal("Failed to register handler: ", err)
	}

	// 运行服务
	if err := service.Run(); err != nil {
		logger.Logger.Fatal("Failed to run service: ", err)
	}
}
