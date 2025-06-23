package handler

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	userProto "github.com/yinxi0607/YixiGroceryAPI/proto/user"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/config"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/model"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/utils"
	"time"
)

type UserHandler struct{}

func (h *UserHandler) Register(ctx context.Context, req *userProto.RegisterRequest, resp *userProto.RegisterResponse) error {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		resp.Code = 500
		resp.Message = "Failed to hash password"
		return err
	}

	user := model.User{
		Username: req.Username,
		Password: hashedPassword,
		Phone:    req.Phone,
	}
	if err := config.DB.Create(&user).Error; err != nil {
		resp.Code = 400
		resp.Message = "Username already exists"
		return err
	}

	resp.Code = 0
	resp.Message = "Success"
	resp.Data = &userProto.User{
		Id:       uint32(user.ID),
		Username: user.Username,
		Phone:    user.Phone,
		Address:  user.Address,
		Points:   int32(user.Points),
	}
	return nil
}

func (h *UserHandler) Login(ctx context.Context, req *userProto.LoginRequest, resp *userProto.LoginResponse) error {
	var user model.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		resp.Code = 400
		resp.Message = "User not found"
		return err
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		resp.Code = 400
		resp.Message = "Invalid password"
		return nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenStr, err := token.SignedString([]byte("your_jwt_secret"))
	if err != nil {
		resp.Code = 500
		resp.Message = "Failed to generate token"
		return err
	}

	resp.Code = 0
	resp.Message = "Success"
	resp.Token = tokenStr
	return nil
}

func (h *UserHandler) GetUserInfo(ctx context.Context, req *userProto.GetUserInfoRequest, resp *userProto.GetUserInfoResponse) error {
	var user model.User
	if err := config.DB.First(&user, req.UserId).Error; err != nil {
		resp.Code = 404
		resp.Message = "User not found"
		return err
	}

	resp.Code = 0
	resp.Message = "Success"
	resp.Data = &userProto.User{
		Id:       uint32(user.ID),
		Username: user.Username,
		Phone:    user.Phone,
		Address:  user.Address,
		Points:   int32(user.Points),
	}
	return nil
}

func (h *UserHandler) AddAddress(ctx context.Context, req *userProto.AddAddressRequest, resp *userProto.AddAddressResponse) error {
	// 如果设置默认地址，清除其他默认地址
	if req.IsDefault {
		config.DB.Model(&model.Address{}).Where("user_id = ? AND is_default = ?", req.UserId, true).Update("is_default", false)
	}

	address := model.Address{
		UserID:        uint(req.UserId),
		ReceiverName:  req.ReceiverName,
		Phone:         req.Phone,
		AddressDetail: req.AddressDetail,
		IsDefault:     req.IsDefault,
	}
	if err := config.DB.Create(&address).Error; err != nil {
		resp.Code = 500
		resp.Message = "Failed to add address"
		return err
	}

	resp.Code = 0
	resp.Message = "Success"
	resp.Data = &userProto.Address{
		Id:            uint32(address.ID),
		UserId:        uint32(address.UserID),
		ReceiverName:  address.ReceiverName,
		Phone:         address.Phone,
		AddressDetail: address.AddressDetail,
		IsDefault:     address.IsDefault,
	}
	return nil
}

func (h *UserHandler) UpdateAddress(ctx context.Context, req *userProto.UpdateAddressRequest, resp *userProto.UpdateAddressResponse) error {
	var address model.Address
	if err := config.DB.Where("id = ? AND user_id = ?", req.Id, req.UserId).First(&address).Error; err != nil {
		resp.Code = 404
		resp.Message = "Address not found"
		return err
	}

	// 如果设置默认地址，清除其他默认地址
	if req.IsDefault {
		config.DB.Model(&model.Address{}).Where("user_id = ? AND is_default = ?", req.UserId, true).Update("is_default", false)
	}

	address.ReceiverName = req.ReceiverName
	address.Phone = req.Phone
	address.AddressDetail = req.AddressDetail
	address.IsDefault = req.IsDefault
	if err := config.DB.Save(&address).Error; err != nil {
		resp.Code = 500
		resp.Message = "Failed to update address"
		return err
	}

	resp.Code = 0
	resp.Message = "Success"
	resp.Data = &userProto.Address{
		Id:            uint32(address.ID),
		UserId:        uint32(address.UserID),
		ReceiverName:  address.ReceiverName,
		Phone:         address.Phone,
		AddressDetail: address.AddressDetail,
		IsDefault:     address.IsDefault,
	}
	return nil
}

func (h *UserHandler) DeleteAddress(ctx context.Context, req *userProto.DeleteAddressRequest, resp *userProto.DeleteAddressResponse) error {
	if err := config.DB.Where("id = ? AND user_id = ?", req.Id, req.UserId).Delete(&model.Address{}).Error; err != nil {
		resp.Code = 404
		resp.Message = "Address not found"
		return err
	}

	resp.Code = 0
	resp.Message = "Success"
	return nil
}

func (h *UserHandler) GetAddresses(ctx context.Context, req *userProto.GetAddressesRequest, resp *userProto.GetAddressesResponse) error {
	var addresses []model.Address
	config.DB.Where("user_id = ?", req.UserId).Find(&addresses)

	resp.Code = 0
	resp.Message = "Success"
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
	return nil
}
