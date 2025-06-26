package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client
var Ctx = context.Background()

func InitDB() {
	// Load .env file in development (optional)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Retrieve MySQL configuration from environment variables
	mysqlUser := getEnv("MYSQL_USER", "user")
	mysqlPassword := getEnv("MYSQL_PASSWORD", "password")
	mysqlHost := getEnv("MYSQL_HOST", "yixi-user-service-db")
	mysqlPort := getEnv("MYSQL_PORT", "3306")
	mysqlDB := getEnv("MYSQL_DATABASE", "user_service_db")

	// Construct MySQL DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDB)

	// Connect to MySQL
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	err = DB.AutoMigrate(&model.User{}, &model.Address{})
	if err != nil {
		return
	}
	log.Println("Connected to MySQL and migrated models")

	// Retrieve Redis configuration from environment variables
	redisAddr := getEnv("REDIS_ADDR", "redis:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := getEnvAsInt("REDIS_DB", 0)

	// Connect to Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
