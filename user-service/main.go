package main

import (
	"github.com/yinxi0607/YixiGroceryAPI/user-service/handler"
	pb "github.com/yinxi0607/YixiGroceryAPI/user-service/proto"

	"go-micro.dev/v5"
)

func main() {
	// Create service
	service := micro.New("user-service")

	// Initialize service
	service.Init()

	// Register handler
	err := pb.RegisterUserServiceHandler(service.Server(), handler.New())
	if err != nil {
		return
	}

	// Run service
	err = service.Run()
	if err != nil {
		return
	}
}
