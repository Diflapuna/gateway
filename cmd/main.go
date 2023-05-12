package main

import "gateway/internal/service"

func main() {
	gateway := service.NewGateway()
	err := gateway.Start()
	if err != nil {
		gateway.Logger.Fatal("Can't start gateway: ", err)
	}
}
