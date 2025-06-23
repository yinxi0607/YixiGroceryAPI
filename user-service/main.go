package main

import (
	"log"

	userProto "github.com/yinxi0607/YixiGroceryAPI/proto/user"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/config"
	"github.com/yinxi0607/YixiGroceryAPI/user-service/handler"
	"go-micro.dev/v5"
)

func main() {
	// 创建 Go-Micro 服务
	service := micro.NewService(
		micro.Name("go.micro.srv.user.yixi"),
		micro.Address(":8081"),
	)

	// 初始化
	config.InitDB()
	userProto.RegisterUserServiceHandler(service.Server(), &handler.UserHandler{})

	// 启动服务
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
