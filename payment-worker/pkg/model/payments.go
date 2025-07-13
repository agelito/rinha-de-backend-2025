package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Payment struct {
	CorrelationId string          `json:"correlationId"`
	Amount        decimal.Decimal `json:"amount"`
	RequestedAt   time.Time       `json:"requestedAt" time_format:"2006-01-02T15:04:05Z07:00"`
}
