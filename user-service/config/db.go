package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB
var RedisClient *redis.Client
var Ctx = context.Background()

func InitDB() {
	dsn := "user:password@tcp(yixi-user-service-db:3306)/user_service_db?charset=utf8mb4&parseTime=True"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL")
	}
	DB.AutoMigrate(&model.User{}, &model.Address{})

	RedisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis")
	}
}
