package subjects

import "fmt"

const (
	SubjectPaymentsProcess    = "payments.process"
	subjectPaymentsConfirmFmt = "payments.confirm.%s"
)

func NewPaymentsConfirmChannel(correlationId string) string {
	return fmt.Sprintf(subjectPaymentsConfirmFmt, correlationId)
}
