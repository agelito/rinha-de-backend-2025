package handler

import (
	"context"
	"fmt"
	"time"

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

	confirmSubject := subjects.NewPaymentsConfirmChannel(payment.CorrelationId.String())

	ch := make(chan *nats.Msg)
	sub, err := h.nc.ChanSubscribe(confirmSubject, ch)

	if err != nil {
		return err
	}

	defer sub.Unsubscribe()

	if err := h.nc.Publish(subjects.SubjectPaymentsProcess, msgBytes); err != nil {
		return err
	}

	if !h.waitForChannelMessage(ch) {
		return fmt.Errorf("timed out waiting for `%v` to process", payment.CorrelationId)
	}

	return nil
}

func (h *PaymentsHandler) waitForChannelMessage(ch chan *nats.Msg) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-ch:
		return true
	case <-ctx.Done():
		return false
	}
}
