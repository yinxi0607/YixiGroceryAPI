package handler

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	userProto "github.com/yinxi0607/YixiGroceryAPI/proto/user"
	"google.golang.org/grpc"
)

type UserHandler struct {
	client userProto.UserServiceClient
}

func NewUserHandler(conn *grpc.ClientConn) *UserHandler {
	return &UserHandler{client: userProto.NewUserServiceClient(conn)}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req userProto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	resp, err := h.client.Register(context.Background(), &req)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req userProto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	resp, err := h.client.Login(context.Background(), &req)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	resp, err := h.client.GetUserInfo(context.Background(), &userProto.GetUserInfoRequest{
		UserId: uint32(userID.(uint)),
	})
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *UserHandler) AddAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		ReceiverName  string `json:"receiver_name"`
		Phone         string `json:"phone"`
		AddressDetail string `json:"address_detail"`
		IsDefault     bool   `json:"is_default"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	resp, err := h.client.AddAddress(context.Background(), &userProto.AddAddressRequest{
		UserId:        uint32(userID.(uint)),
		ReceiverName:  req.ReceiverName,
		Phone:         req.Phone,
		AddressDetail: req.AddressDetail,
		IsDefault:     req.IsDefault,
	})
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *UserHandler) UpdateAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	var req struct {
		ReceiverName  string `json:"receiver_name"`
		Phone         string `json:"phone"`
		AddressDetail string `json:"address_detail"`
		IsDefault     bool   `json:"is_default"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	resp, err := h.client.UpdateAddress(context.Background(), &userProto.UpdateAddressRequest{
		Id:            uint32(atoi(id)),
		UserId:        uint32(userID.(uint)),
		ReceiverName:  req.ReceiverName,
		Phone:         req.Phone,
		AddressDetail: req.AddressDetail,
		IsDefault:     req.IsDefault,
	})
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *UserHandler) DeleteAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	resp, err := h.client.DeleteAddress(context.Background(), &userProto.DeleteAddressRequest{
		Id:     uint32(atoi(id)),
		UserId: uint32(userID.(uint)),
	})
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *UserHandler) GetAddresses(c *gin.Context) {
	userID, _ := c.Get("user_id")
	resp, err := h.client.GetAddresses(context.Background(), &userProto.GetAddressesRequest{
		UserId: uint32(userID.(uint)),
	})
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
