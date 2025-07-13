package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/agelito/rinha-de-backend-2025/api/pkg/model"
	"github.com/agelito/rinha-de-backend-2025/messages/model/payments"
	pb "github.com/agelito/rinha-de-backend-2025/messages/model/payments"
	"github.com/agelito/rinha-de-backend-2025/messages/subjects"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

var (
	PaymentFailed        = fmt.Errorf("payment processing failed")
	PaymentTimeout       = fmt.Errorf("payment timed out")
	PaymentInternalError = fmt.Errorf("internal server error")
)

type PaymentsHandler struct {
	nc *nats.Conn
}

func NewPaymentsHandler(nc *nats.Conn) *PaymentsHandler {
	return &PaymentsHandler{nc: nc}
}

func (h *PaymentsHandler) Payment(payment *model.Payment) error {
	requestedAt := time.Now().UTC()

	msg := &pb.Payment{
		CorrelationId: payment.CorrelationId.String(),
		Amount:        payment.Amount.String(),
		RequestedAt:   requestedAt.Format(time.RFC3339),
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

	return h.waitForPaymentResult(ch)
}

func (h *PaymentsHandler) waitForPaymentResult(ch chan *nats.Msg) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case msg := <-ch:
		var result payments.PaymentResult
		if err := proto.Unmarshal(msg.Data, &result); err != nil {
			return PaymentInternalError
		}

		if !result.Successful {
			return PaymentFailed
		}

		return nil
	case <-ctx.Done():
		return PaymentTimeout
	}
}
