package main

import (
	"log"

	"github.com/agelito/rinha-de-backend-2025/api/pkg/handler"
	"github.com/agelito/rinha-de-backend-2025/api/pkg/service"
)

func main() {
	payments := handler.NewPaymentsHandler()
	httpService := service.NewHttpService(payments)

	if err := httpService.Run(":3001"); err != nil {
		log.Fatal(err)
	}
}
