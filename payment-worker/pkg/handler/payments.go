package handler

import (
	"fmt"

	"github.com/agelito/rinha-de-backend-2025/messages/subjects"
	"github.com/nats-io/nats.go"
)

type PaymentsHandler struct {
	nc *nats.Conn
}

func NewPaymentsHandler(nc *nats.Conn) *PaymentsHandler {
	return &PaymentsHandler{nc: nc}
}

func (h *PaymentsHandler) ProcessPayment(correlationId string, amount string) error {
	fmt.Printf("payment: correlationId: %v, amount: %v\n", correlationId, amount)

	confirmSubject := subjects.NewPaymentsConfirmChannel(correlationId)
	if err := h.nc.Publish(confirmSubject, []byte{}); err != nil {
		return err
	}

	return nil
}
