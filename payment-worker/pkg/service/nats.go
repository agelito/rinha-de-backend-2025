package service

import (
	"context"
	"log"

	pb "github.com/agelito/rinha-de-backend-2025/messages/model/payments"
	"github.com/agelito/rinha-de-backend-2025/messages/subjects"
	"github.com/agelito/rinha-de-backend-2025/payment-worker/pkg/handler"
	worker "github.com/agelito/rinha-de-backend-2025/payment-worker/pkg/worker"
	"github.com/nats-io/nats.go"
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
					log.Printf("error deserializing message: %v\n", err)
				}

				if err := s.handler.ProcessPayment(paymentMsg.CorrelationId, paymentMsg.Amount); err != nil {
					log.Printf("error processing payment: %v\n", err)
					// TODO: retry the payment (add to retry queue)
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
