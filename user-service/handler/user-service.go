package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/yinxi0607/YixiGroceryAPI/logger"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/model"
	pb "github.com/yinxi0607/YixiGroceryAPI/user-service/proto"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"time"
)

type UserService struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	pb.UserService
}

// New creates a new UserService instance, injecting GORM database and Redis client
func New(db *gorm.DB, redisClient *redis.Client) *UserService {
	return &UserService{
		DB:          db,
		RedisClient: redisClient,
	}
}

// generateUUID generates a UUID
func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// generateJWT generates a JWT access token
func generateJWT(userID string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // 访问令牌 1 小时有效
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(jwtSecret))
}

// generateRefreshToken generates a refresh token
func generateRefreshToken(userID string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 168).Unix(), // 刷新令牌 7 天有效
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	})

	return token.SignedString([]byte(jwtSecret))
}

// Register implements user registration
func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest, rsp *pb.RegisterResponse) error {
	logger.Logger.WithFields(logrus.Fields{
		"username": req.Username,
	}).Info("Register request received")

	if req.Username == "" || req.Password == "" || req.Email == "" {
		logger.Logger.Error("Invalid input: username, password, or email is empty")
		rsp.Success = false
		rsp.Message = "Username, password, and email are required"
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Error("Failed to hash password: ", err)
		rsp.Success = false
		rsp.Message = "Failed to hash password"
		return err
	}

	user := model.User{
		ID:       generateUUID(),
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Address:  req.Address,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		logger.Logger.Error("Failed to create user: ", err)
		rsp.Success = false
		rsp.Message = "Failed to create user"
		return err
	}

	rsp.UserId = user.ID
	rsp.Success = true
	rsp.Message = "User registered successfully"
	return nil
}

// Login implements user login with JWT and refresh token
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	logger.Logger.WithFields(logrus.Fields{
		"username": req.Username,
	}).Info("Login request received")

	if req.Username == "" || req.Password == "" {
		logger.Logger.Error("Invalid input: username or password is empty")
		rsp.Success = false
		rsp.Message = "Username and password are required"
		return nil
	}

	var user model.User
	if err := s.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		logger.Logger.Error("User not found: ", err)
		rsp.Success = false
		rsp.Message = "Invalid username or password"
		return nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Logger.Error("Invalid password: ", err)
		rsp.Success = false
		rsp.Message = "Invalid username or password"
		return nil
	}

	// Generate JWT access token
	accessToken, err := generateJWT(user.ID)
	if err != nil {
		logger.Logger.Error("Failed to generate access token: ", err)
		rsp.Success = false
		rsp.Message = "Failed to generate access token"
		return err
	}

	// Generate refresh token
	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		logger.Logger.Error("Failed to generate refresh token: ", err)
		rsp.Success = false
		rsp.Message = "Failed to generate refresh token"
		return err
	}

	// Store refresh token in Redis (7 days TTL)
	ctx = context.Background()
	err = s.RedisClient.Set(ctx, "refresh_token:"+user.ID, refreshToken, 7*24*time.Hour).Err()
	if err != nil {
		logger.Logger.Error("Failed to store refresh token in Redis: ", err)
		rsp.Success = false
		rsp.Message = "Failed to store refresh token"
		return err
	}

	rsp.UserId = user.ID
	rsp.Token = accessToken
	rsp.RefreshToken = refreshToken
	rsp.Success = true
	rsp.Message = "Login successful"
	return nil
}

// RefreshToken implements token refresh
func (s *UserService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest, rsp *pb.RefreshTokenResponse) error {
	logger.Logger.WithFields(logrus.Fields{
		"user_id": req.UserId,
	}).Info("RefreshToken request received")

	if req.UserId == "" || req.RefreshToken == "" {
		logger.Logger.Error("Invalid input: user_id or refresh_token is empty")
		rsp.Success = false
		rsp.Message = "User ID and refresh token are required"
		return nil
	}

	// Verify refresh token from Redis
	storedToken, err := s.RedisClient.Get(ctx, "refresh_token:"+req.UserId).Result()
	if err == redis.Nil || storedToken != req.RefreshToken {
		logger.Logger.Error("Invalid or expired refresh token")
		rsp.Success = false
		rsp.Message = "Invalid or expired refresh token"
		return nil
	}

	// Generate new access token
	newAccessToken, err := generateJWT(req.UserId)
	if err != nil {
		logger.Logger.Error("Failed to generate new access token: ", err)
		rsp.Success = false
		rsp.Message = "Failed to generate new access token"
		return err
	}

	rsp.Token = newAccessToken
	rsp.Success = true
	rsp.Message = "Token refreshed successfully"
	return nil
}

// GetUser retrieves user information with caching
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest, rsp *pb.GetUserResponse) error {
	logger.Logger.WithFields(logrus.Fields{
		"user_id": req.UserId,
	}).Info("GetUser request received")

	// Check Redis cache
	cacheKey := "user:" + req.UserId
	cachedUser, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit
		var user model.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err != nil {
			logger.Logger.Error("Failed to unmarshal cached user: ", err)
		} else {
			logger.Logger.Info("Cache hit for user: ", req.UserId)
			rsp.UserId = user.ID
			rsp.Username = user.Username
			rsp.Email = user.Email
			rsp.Address = user.Address
			rsp.Success = true
			rsp.Message = "User retrieved from cache"
			return nil
		}
	}

	// Cache miss, query database
	var user model.User
	if err := s.DB.Where("id = ?", req.UserId).First(&user).Error; err != nil {
		logger.Logger.Error("User not found: ", err)
		rsp.Success = false
		rsp.Message = "User not found"
		return nil
	}

	// Cache user data in Redis (1 hour TTL)
	userData, err := json.Marshal(user)
	if err != nil {
		logger.Logger.Error("Failed to marshal user for cache: ", err)
	} else {
		err = s.RedisClient.Set(ctx, cacheKey, userData, time.Hour).Err()
		if err != nil {
			logger.Logger.Error("Failed to cache user: ", err)
		}
	}

	rsp.UserId = user.ID
	rsp.Username = user.Username
	rsp.Email = user.Email
	rsp.Address = user.Address
	rsp.Success = true
	rsp.Message = "User retrieved successfully"
	return nil
}

// UpdateUser updates user information and invalidates cache
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest, rsp *pb.UpdateUserResponse) error {
	logger.Logger.WithFields(logrus.Fields{
		"user_id": req.UserId,
	}).Info("UpdateUser request received")

	var user model.User
	if err := s.DB.Where("id = ?", req.UserId).First(&user).Error; err != nil {
		logger.Logger.Error("User not found: ", err)
		rsp.Success = false
		rsp.Message = "User not found"
		return nil
	}

	// Update fields
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Logger.Error("Failed to hash password: ", err)
			rsp.Success = false
			rsp.Message = "Failed to hash password"
			return err
		}
		user.Password = string(hashedPassword)
	}

	if err := s.DB.Save(&user).Error; err != nil {
		logger.Logger.Error("Failed to update user: ", err)
		rsp.Success = false
		rsp.Message = "Failed to update user"
		return err
	}

	// Invalidate cache
	cacheKey := "user:" + req.UserId
	if err := s.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
		logger.Logger.Error("Failed to invalidate cache: ", err)
	}

	rsp.Success = true
	rsp.Message = "User updated successfully"
	return nil
}

// DeleteUser deletes a user and invalidates cache
func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest, rsp *pb.DeleteUserResponse) error {
	logger.Logger.WithFields(logrus.Fields{
		"user_id": req.UserId,
	}).Info("DeleteUser request received")

	if err := s.DB.Where("id = ?", req.UserId).Delete(&model.User{}).Error; err != nil {
		logger.Logger.Error("Failed to delete user: ", err)
		rsp.Success = false
		rsp.Message = "Failed to delete user"
		return err
	}

	// Invalidate cache
	cacheKey := "user:" + req.UserId
	if err := s.RedisClient.Del(ctx, cacheKey).Err(); err != nil {
		logger.Logger.Error("Failed to invalidate cache: ", err)
	}

	// Delete refresh token
	if err := s.RedisClient.Del(ctx, "refresh_token:"+req.UserId).Err(); err != nil {
		logger.Logger.Error("Failed to delete refresh token: ", err)
	}

	rsp.Success = true
	rsp.Message = "User deleted successfully"
	return nil
}

// ListUsers lists all users with pagination
func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest, rsp *pb.ListUsersResponse) error {
	logger.Logger.WithFields(logrus.Fields{
		"page":      req.Page,
		"page_size": req.PageSize,
	}).Info("ListUsers request received")

	var users []model.User
	offset := (req.Page - 1) * req.PageSize
	if err := s.DB.Limit(int(req.PageSize)).Offset(int(offset)).Find(&users).Error; err != nil {
		logger.Logger.Error("Failed to list users: ", err)
		rsp.Success = false
		rsp.Message = "Failed to list users"
		return err
	}

	rsp.Users = make([]*pb.User, len(users))
	for i, user := range users {
		rsp.Users[i] = &pb.User{
			UserId:   user.ID,
			Username: user.Username,
			Email:    user.Email,
			Address:  user.Address,
		}
	}

	rsp.Success = true
	rsp.Message = "Users retrieved successfully"
	return nil
}
