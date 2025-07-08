package handler

import (
	"github.com/agelito/rinha-de-backend-2025/api/pkg/model"
)

type PaymentsHandler struct{}

func NewPaymentsHandler() *PaymentsHandler {
	return &PaymentsHandler{}
}

func (h *PaymentsHandler) Payment(payment *model.Payment) error {
	return nil
}
