package main

import (
	"os"

	"github.com/agelito/rinha-de-backend-2025/api/pkg/handler"
	"github.com/agelito/rinha-de-backend-2025/api/pkg/service"
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

	payments := handler.NewPaymentsHandler(nc)
	httpService := service.NewHttpService(payments)

	if err := httpService.Run(":3001"); err != nil {
		log.Fatal(err)
	}
}
