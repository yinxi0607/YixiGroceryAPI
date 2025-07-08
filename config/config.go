package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/yinxi0607/YixiGroceryAPI/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

// DatabaseConfig 存储数据库配置
type DatabaseConfig struct {
	Driver          string        // 数据库驱动：postgres, mysql 等
	DSN             string        // 数据库连接字符串
	Host            string        // 主机
	User            string        // 用户名
	Password        string        // 密码
	DBName          string        // 数据库名
	Port            string        // 端口
	SSLMode         string        // SSL 模式（仅限 PostgreSQL）
	MaxIdleConns    int           // 最大空闲连接数
	MaxOpenConns    int           // 最大打开连接数
	ConnMaxLifetime time.Duration // 连接最大生命周期
}

// LoadDBConfig 从环境变量加载数据库配置
func LoadDBConfig(serviceName string) DatabaseConfig {
	cfg := DatabaseConfig{
		Driver:   os.Getenv(fmt.Sprintf("%s_DB_DRIVER", serviceName)),
		Host:     os.Getenv(fmt.Sprintf("%s_DB_HOST", serviceName)),
		User:     os.Getenv(fmt.Sprintf("%s_DB_USER", serviceName)),
		Password: os.Getenv(fmt.Sprintf("%s_DB_PASSWORD", serviceName)),
		DBName:   os.Getenv(fmt.Sprintf("%s_DB_NAME", serviceName)),
		Port:     os.Getenv(fmt.Sprintf("%s_DB_PORT", serviceName)),
		SSLMode:  os.Getenv(fmt.Sprintf("%s_DB_SSLMODE", serviceName)),
	}

	// 加载连接池配置
	maxIdleConns, _ := strconv.Atoi(os.Getenv(fmt.Sprintf("%s_DB_MAX_IDLE_CONNS", serviceName)))
	if maxIdleConns == 0 {
		maxIdleConns = 10 // 默认值
	}
	maxOpenConns, _ := strconv.Atoi(os.Getenv(fmt.Sprintf("%s_DB_MAX_OPEN_CONNS", serviceName)))
	if maxOpenConns == 0 {
		maxOpenConns = 100 // 默认值
	}
	connMaxLifetime, _ := strconv.Atoi(os.Getenv(fmt.Sprintf("%s_DB_CONN_MAX_LIFETIME", serviceName)))
	if connMaxLifetime == 0 {
		connMaxLifetime = 3600 // 默认 1 小时（秒）
	}

	cfg.MaxIdleConns = maxIdleConns
	cfg.MaxOpenConns = maxOpenConns
	cfg.ConnMaxLifetime = time.Duration(connMaxLifetime) * time.Second

	// 如果提供了特定服务的 DATABASE_URL，直接使用
	if dsn := os.Getenv(fmt.Sprintf("%s_DATABASE_URL", serviceName)); dsn != "" {
		cfg.DSN = dsn
		return cfg
	}

	// 如果没有 DATABASE_URL，使用单独的字段构建 DSN
	if cfg.Driver == "" {
		cfg.Driver = "postgres"
	}
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.User == "" {
		cfg.User = "postgres"
	}
	if cfg.Password == "" {
		cfg.Password = "yourpassword"
	}
	if cfg.DBName == "" {
		cfg.DBName = fmt.Sprintf("%s_db", serviceName)
	}
	if cfg.Port == "" {
		cfg.Port = "5432"
	}
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	// 根据驱动构建 DSN
	switch cfg.Driver {
	case "postgres":
		cfg.DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
	case "mysql":
		cfg.DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	default:
		logger.Logger.Warnf("Unsupported database driver: %s, defaulting to postgres", cfg.Driver)
		cfg.DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
	}

	return cfg
}

// InitDB 初始化 GORM 数据库连接
func InitDB(cfg DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch cfg.Driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	default:
		logger.Logger.Errorf("Unsupported database driver: %s", cfg.Driver)
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	if err != nil {
		logger.Logger.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		logger.Logger.Errorf("Failed to get sql.DB: %v", err)
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	logger.Logger.Info("Database connected with connection pool")
	return db, nil
}

// RedisConfig 存储 Redis 配置
type RedisConfig struct {
	Addr     string // Redis 地址（host:port）
	Password string // Redis 密码
	DB       int    // Redis 数据库编号
}

// LoadRedisConfig 从环境变量加载 Redis 配置
func LoadRedisConfig(serviceName string) RedisConfig {
	cfg := RedisConfig{
		Addr:     os.Getenv(fmt.Sprintf("%s_REDIS_ADDR", serviceName)),
		Password: os.Getenv(fmt.Sprintf("%s_REDIS_PASSWORD", serviceName)),
	}

	db, _ := strconv.Atoi(os.Getenv(fmt.Sprintf("%s_REDIS_DB", serviceName)))
	cfg.DB = db

	if cfg.Addr == "" {
		cfg.Addr = "localhost:6379" // 默认值
	}

	return cfg
}

// InitRedis 初始化 Redis 客户端
func InitRedis(cfg RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试 Redis 连接
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Logger.Errorf("Failed to connect to Redis: %v", err)
		return nil, err
	}

	logger.Logger.Info("Redis connected")
	return client, nil
}
