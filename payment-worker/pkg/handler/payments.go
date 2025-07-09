package handler

import "fmt"

type PaymentsHandler struct{}

func NewPaymentsHandler() *PaymentsHandler {
	return &PaymentsHandler{}
}

func (h *PaymentsHandler) ProcessPayment(correlationId string, amount string) error {
	fmt.Printf("payment: correlationId: %v, amount: %v\n", correlationId, amount)
	return nil
}
