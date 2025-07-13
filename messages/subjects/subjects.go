package subjects

import "fmt"

const (
	SubjectServerScore        = "servers.score"
	SubjectServerRanking      = "servers.ranking"
	SubjectPaymentsProcess    = "payments.process"
	subjectPaymentsConfirmFmt = "payments.confirm.%s"
)

func NewPaymentsConfirmChannel(correlationId string) string {
	return fmt.Sprintf(subjectPaymentsConfirmFmt, correlationId)
}
