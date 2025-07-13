package subjects

import "fmt"

const (
	SubjectServerScore        = "servers.score"
	SubjectPaymentsProcess    = "payments.process"
	subjectPaymentsConfirmFmt = "payments.confirm.%s"
)

func NewPaymentsConfirmChannel(correlationId string) string {
	return fmt.Sprintf(subjectPaymentsConfirmFmt, correlationId)
}
