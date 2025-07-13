package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/agelito/rinha-de-backend-2025/messages/model/payments"
	"github.com/agelito/rinha-de-backend-2025/messages/model/servers"
	"github.com/agelito/rinha-de-backend-2025/messages/subjects"
	"github.com/agelito/rinha-de-backend-2025/payment-worker/pkg/model"
	workers "github.com/agelito/rinha-de-backend-2025/payment-worker/pkg/worker"
	"github.com/charmbracelet/log"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type ServerScore struct {
	Server string
	Score  int32
}

type PaymentsHandler struct {
	nc      *nats.Conn
	worker  *workers.Worker
	client  *http.Client
	servers []ServerScore
}

func NewPaymentsHandler(nc *nats.Conn) *PaymentsHandler {
	return &PaymentsHandler{
		nc:     nc,
		worker: workers.NewWorker(),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		servers: make([]ServerScore, 0),
	}
}

func (h *PaymentsHandler) ProcessPayment(payment *model.Payment) error {
	log.Info("process payment", "correlationId", payment.CorrelationId, "amount", payment.Amount, "requestedAt", payment.RequestedAt)

	confirmSubject := subjects.NewPaymentsConfirmChannel(payment.CorrelationId)

	err := h.sendProcessPaymentRequest(payment)

	result := &payments.PaymentResult{
		CorrelationId: payment.CorrelationId,
		Successful:    err == nil,
	}

	if err != nil {
		result.Error = err.Error()
	}

	msgBytes, err := proto.Marshal(result)

	if err != nil {
		return err
	}

	if err := h.nc.Publish(confirmSubject, msgBytes); err != nil {
		return err
	}

	return nil
}

func (h *PaymentsHandler) Run() error {
	ch := make(chan *nats.Msg, 2)
	sub, err := h.nc.ChanSubscribe(subjects.SubjectServerRanking, ch)

	if err != nil {
		return err
	}

	defer sub.Unsubscribe()

	h.worker.Run(h.serverRankingsWorker(ch))
	h.worker.Join()

	return nil
}

func (h *PaymentsHandler) Stop() {
	h.worker.Stop()
	h.worker.Join()
}

func (h *PaymentsHandler) serverRankingsWorker(ch chan *nats.Msg) func(context.Context) {
	return func(workerCtx context.Context) {
		for {
			select {
			case msg := <-ch:
				var ranking servers.ServerRanking
				if err := proto.Unmarshal(msg.Data, &ranking); err != nil {
					log.Error("could not deserialize server rankings", "error", err)
				}

				servers := make([]ServerScore, len(ranking.Servers))

				for idx, srv := range ranking.Servers {
					servers[idx] = ServerScore{
						Server: srv.Server,
						Score:  srv.Score,
					}
				}

				h.servers = servers

			case <-workerCtx.Done():
				return
			}
		}
	}
}

func (h *PaymentsHandler) sendProcessPaymentRequest(payment *model.Payment) error {
	bodyBytes, err := json.Marshal(payment)

	if err != nil {
		return err
	}

	if len(h.servers) == 0 {
		return fmt.Errorf("no servers available")
	}

	// TODO: Handle cases where all servers is in a failing state
	url := fmt.Sprintf("%s/payments", h.servers[0].Server)

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return err
	}

	res, err := h.client.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("could not process payment, unexpected response status: %v", res.StatusCode)
	}

	return nil
}
