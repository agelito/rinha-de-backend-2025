package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Payment struct {
	CorrelationId uuid.UUID       `json:"correlationId"`
	Amount        decimal.Decimal `json:"amount"`
}
