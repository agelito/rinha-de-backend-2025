package service

import (
	"context"
	"time"

	pb "github.com/agelito/rinha-de-backend-2025/messages/model/payments"
	"github.com/agelito/rinha-de-backend-2025/messages/subjects"
	"github.com/agelito/rinha-de-backend-2025/payment-worker/pkg/handler"
	"github.com/agelito/rinha-de-backend-2025/payment-worker/pkg/model"
	worker "github.com/agelito/rinha-de-backend-2025/payment-worker/pkg/worker"
	"github.com/charmbracelet/log"
	"github.com/nats-io/nats.go"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"
)

const (
	queue  = "payment-worker"
	buffer = 128
)

type NatsService struct {
	nc      *nats.Conn
	worker  *worker.Worker
	handler *handler.PaymentsHandler
}

func NewNatsService(nc *nats.Conn, h *handler.PaymentsHandler) *NatsService {
	return &NatsService{
		nc:      nc,
		worker:  worker.NewWorker(),
		handler: h,
	}
}

func (s *NatsService) Run() {
	s.worker.Run(func(workerCtx context.Context) {
		ch := make(chan *nats.Msg, buffer)

		sub, err := s.nc.QueueSubscribeSyncWithChan(subjects.SubjectPaymentsProcess, queue, ch)

		if err != nil {
			log.Fatal(err)
		}

		defer sub.Unsubscribe()

		for {
			select {
			case msg := <-ch:
				var paymentMsg pb.Payment
				if err := proto.Unmarshal(msg.Data, &paymentMsg); err != nil {
					// TODO: add to dlq for error reporting?
					log.Error("could not deserialize message", "error", err)
					continue
				}

				amount, err := decimal.NewFromString(paymentMsg.Amount)

				if err != nil {
					// TODO: add to dlq for error reporting?
					log.Error("could not parse amount", "error", err)
					continue
				}

				requestedAt, err := time.Parse(time.RFC3339, paymentMsg.RequestedAt)

				if err != nil {
					// TODO: add to dlq for error reporting?
					log.Error("could not parse requestedAt", "error", err)
					continue
				}

				payment := &model.Payment{
					CorrelationId: paymentMsg.CorrelationId,
					Amount:        amount,
					RequestedAt:   requestedAt,
				}

				if err := s.handler.ProcessPayment(payment); err != nil {
					log.Error("could not process payment", "error", err)
				}
			case <-workerCtx.Done():
				return
			}
		}
	})

	s.worker.Join()
}

func (s *NatsService) Stop() {
	s.worker.Stop()
}
