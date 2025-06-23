package main

import (
	"github.com/micro/go-micro/v2"
	userProto "github.com/yinxi0607/YixiGroceryAPI/proto/user"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/config"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/handler"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.user.yixi"),
		micro.Address(":8081"),
	)
	service.Init()

	config.InitDB()
	userProto.RegisterUserServiceHandler(service.Server(), &handler.UserHandler{})
	service.Run()
}
