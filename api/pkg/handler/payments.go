package handler

import (
	"github.com/agelito/rinha-de-backend-2025/api/pkg/model"
	pb "github.com/agelito/rinha-de-backend-2025/messages/model/payments"
	"github.com/agelito/rinha-de-backend-2025/messages/subjects"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type PaymentsHandler struct {
	nc *nats.Conn
}

func NewPaymentsHandler(nc *nats.Conn) *PaymentsHandler {
	return &PaymentsHandler{nc: nc}
}

func (h *PaymentsHandler) Payment(payment *model.Payment) error {
	msg := &pb.Payment{
		CorrelationId: payment.CorrelationId.String(),
		Amount:        payment.Amount.String(),
	}

	msgBytes, err := proto.Marshal(msg)

	if err != nil {
		return err
	}

	return h.nc.Publish(subjects.SubjectPaymentsProcess, msgBytes)
}
