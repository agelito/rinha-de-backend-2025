package errors

import "fmt"

var (
	PaymentsUnsuccessful = fmt.Errorf("unsuccessful payment")
)
