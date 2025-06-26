package handler

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	userProto "github.com/yinxi0607/YixiGroceryAPI/proto/user"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/config"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/model"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/utils"
	"time"
)

type UserHandler struct {
	userProto.UnimplementedUserServiceServer
}

func (h *UserHandler) Register(ctx context.Context, req *userProto.RegisterRequest) (*userProto.RegisterResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return &userProto.RegisterResponse{Code: 500, Message: "Failed to hash password"}, err
	}

	user := model.User{
		Username: req.Username,
		Password: hashedPassword,
		Phone:    req.Phone,
	}
	if err = config.DB.WithContext(ctx).Create(&user).Error; err != nil {
		return &userProto.RegisterResponse{Code: 400, Message: "Username already exists"}, err
	}

	return &userProto.RegisterResponse{
		Code:    0,
		Message: "Success",
		Data: &userProto.User{
			Id:       uint32(user.ID),
			Username: user.Username,
			Phone:    user.Phone,
			Address:  user.Address,
			Points:   int32(user.Points),
		},
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *userProto.LoginRequest) (*userProto.LoginResponse, error) {
	var user model.User
	if err := config.DB.WithContext(ctx).Where("username = ?", req.Username).First(&user).Error; err != nil {
		return &userProto.LoginResponse{Code: 400, Message: "User not found"}, err
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return &userProto.LoginResponse{Code: 400, Message: "Invalid password"}, nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenStr, err := token.SignedString([]byte("your_jwt_secret"))
	if err != nil {
		return &userProto.LoginResponse{Code: 500, Message: "Failed to generate token"}, err
	}

	return &userProto.LoginResponse{
		Code:    0,
		Message: "Success",
		Token:   tokenStr,
	}, nil
}

func (h *UserHandler) GetUserInfo(ctx context.Context, req *userProto.GetUserInfoRequest) (*userProto.GetUserInfoResponse, error) {
	var user model.User
	if err := config.DB.WithContext(ctx).First(&user, req.UserId).Error; err != nil {
		return &userProto.GetUserInfoResponse{Code: 404, Message: "User not found"}, err
	}

	return &userProto.GetUserInfoResponse{
		Code:    0,
		Message: "Success",
		Data: &userProto.User{
			Id:       uint32(user.ID),
			Username: user.Username,
			Phone:    user.Phone,
			Address:  user.Address,
			Points:   int32(user.Points),
		},
	}, nil
}

func (h *UserHandler) AddAddress(ctx context.Context, req *userProto.AddAddressRequest) (*userProto.AddAddressResponse, error) {
	if req.IsDefault {
		config.DB.WithContext(ctx).Model(&model.Address{}).Where("user_id = ? AND is_default = ?", req.UserId, true).Update("is_default", false)
	}

	address := model.Address{
		UserID:        uint(req.UserId),
		ReceiverName:  req.ReceiverName,
		Phone:         req.Phone,
		AddressDetail: req.AddressDetail,
		IsDefault:     req.IsDefault,
	}
	if err := config.DB.Create(&address).Error; err != nil {
		return &userProto.AddAddressResponse{Code: 500, Message: "Failed to add address"}, err
	}

	return &userProto.AddAddressResponse{
		Code:    0,
		Message: "Success",
		Data: &userProto.Address{
			Id:            uint32(address.ID),
			UserId:        uint32(address.UserID),
			ReceiverName:  address.ReceiverName,
			Phone:         address.Phone,
			AddressDetail: address.AddressDetail,
			IsDefault:     address.IsDefault,
		},
	}, nil
}

func (h *UserHandler) UpdateAddress(ctx context.Context, req *userProto.UpdateAddressRequest) (*userProto.UpdateAddressResponse, error) {
	var address model.Address
	if err := config.DB.WithContext(ctx).Where("id = ? AND user_id = ?", req.Id, req.UserId).First(&address).Error; err != nil {
		return &userProto.UpdateAddressResponse{Code: 404, Message: "Address not found"}, err
	}

	if req.IsDefault {
		config.DB.Model(&model.Address{}).Where("user_id = ? AND is_default = ?", req.UserId, true).Update("is_default", false)
	}

	address.ReceiverName = req.ReceiverName
	address.Phone = req.Phone
	address.AddressDetail = req.AddressDetail
	address.IsDefault = req.IsDefault
	if err := config.DB.Save(&address).Error; err != nil {
		return &userProto.UpdateAddressResponse{Code: 500, Message: "Failed to update address"}, err
	}

	return &userProto.UpdateAddressResponse{
		Code:    0,
		Message: "Success",
		Data: &userProto.Address{
			Id:            uint32(address.ID),
			UserId:        uint32(address.UserID),
			ReceiverName:  address.ReceiverName,
			Phone:         address.Phone,
			AddressDetail: address.AddressDetail,
			IsDefault:     address.IsDefault,
		},
	}, nil
}

func (h *UserHandler) DeleteAddress(ctx context.Context, req *userProto.DeleteAddressRequest) (*userProto.DeleteAddressResponse, error) {
	if err := config.DB.WithContext(ctx).Where("id = ? AND user_id = ?", req.Id, req.UserId).Delete(&model.Address{}).Error; err != nil {
		return &userProto.DeleteAddressResponse{Code: 404, Message: "Address not found"}, err
	}

	return &userProto.DeleteAddressResponse{
		Code:    0,
		Message: "Success",
	}, nil
}

func (h *UserHandler) GetAddresses(ctx context.Context, req *userProto.GetAddressesRequest) (*userProto.GetAddressesResponse, error) {
	var addresses []model.Address
	config.DB.WithContext(ctx).Where("user_id = ?", req.UserId).Find(&addresses)

	resp := &userProto.GetAddressesResponse{
		Code:    0,
		Message: "Success",
	}
	for _, addr := range addresses {
		resp.Addresses = append(resp.Addresses, &userProto.Address{
			Id:            uint32(addr.ID),
			UserId:        uint32(addr.UserID),
			ReceiverName:  addr.ReceiverName,
			Phone:         addr.Phone,
			AddressDetail: addr.AddressDetail,
			IsDefault:     addr.IsDefault,
		})
	}
	return resp, nil
}
