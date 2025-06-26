package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"github.com/yinxi0607/YixiGroceryAPI/logger"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/model"
	pb "github.com/yinxi0607/YixiGroceryAPI/user-service/proto"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
	pb.UserService
}

// New creates a new UserService instance, injecting GORM database
func New(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// generateUUID generates a UUID
func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Register implements user registration
func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest, rsp *pb.RegisterResponse) error {
	logger.Logger.WithFields(logrus.Fields{
		"username": req.Username,
	}).Info("Register request received")

	// Validate input
	if req.Username == "" || req.Password == "" || req.Email == "" {
		logger.Logger.Error("Invalid input: username, password, or email is empty")
		rsp.Success = false
		rsp.Message = "Username, password, and email are required"
		return nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Error("Failed to hash password: ", err)
		rsp.Success = false
		rsp.Message = "Failed to hash password"
		return err
	}

	// Create user
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

// Login implements user login
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	logger.Logger.WithFields(logrus.Fields{
		"username": req.Username,
	}).Info("Login request received")

	// Validate input
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

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Logger.Error("Invalid password: ", err)
		rsp.Success = false
		rsp.Message = "Invalid username or password"
		return nil
	}

	// TODO: Implement JWT token generation
	rsp.UserId = user.ID
	rsp.Token = "jwt_token_placeholder" // Replace with actual JWT implementation
	rsp.Success = true
	rsp.Message = "Login successful"
	return nil
}

// GetUser retrieves user information
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest, rsp *pb.GetUserResponse) error {
	logger.Logger.WithFields(logrus.Fields{
		"user_id": req.UserId,
	}).Info("GetUser request received")

	var user model.User
	if err := s.DB.Where("id = ?", req.UserId).First(&user).Error; err != nil {
		logger.Logger.Error("User not found: ", err)
		rsp.Success = false
		rsp.Message = "User not found"
		return nil
	}

	rsp.UserId = user.ID
	rsp.Username = user.Username
	rsp.Email = user.Email
	rsp.Address = user.Address
	rsp.Success = true
	rsp.Message = "User retrieved successfully"
	return nil
}

// UpdateUser updates user information
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

	rsp.Success = true
	rsp.Message = "User updated successfully"
	return nil
}

// DeleteUser deletes a user
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
