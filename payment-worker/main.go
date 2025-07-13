package main

import (
	"log"

	"github.com/agelito/rinha-de-backend-2025/payment-worker/pkg/handler"
	"github.com/agelito/rinha-de-backend-2025/payment-worker/pkg/service"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	handler := handler.NewPaymentsHandler(nc)
	service := service.NewNatsService(nc, handler)

	service.Run()
}
