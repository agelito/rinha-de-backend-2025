package main

import (
	"os"

	"github.com/agelito/rinha-de-backend-2025/payment-worker/pkg/handler"
	"github.com/agelito/rinha-de-backend-2025/payment-worker/pkg/service"
	"github.com/charmbracelet/log"
	"github.com/nats-io/nats.go"
)

func main() {
	natsUrl := nats.DefaultURL

	if envUrl := os.Getenv("NATS_URL"); envUrl != "" {
		natsUrl = envUrl
	}

	log.Info("connecting to nats", "url", natsUrl)

	nc, err := nats.Connect(natsUrl)

	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	handler := handler.NewPaymentsHandler(nc)
	service := service.NewNatsService(nc, handler)

	service.Run()
}
