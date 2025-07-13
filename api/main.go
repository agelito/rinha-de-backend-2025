package main

import (
	"log"

	"github.com/agelito/rinha-de-backend-2025/api/pkg/handler"
	"github.com/agelito/rinha-de-backend-2025/api/pkg/service"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	payments := handler.NewPaymentsHandler(nc)
	httpService := service.NewHttpService(payments)

	if err := httpService.Run(":9999"); err != nil {
		log.Fatal(err)
	}
}
