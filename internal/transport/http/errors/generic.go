package errors

import "fmt"

type GenericError struct {
	Code     string
	Reason   string
	HttpCode int
	Message  string
}

func (e *GenericError) Error() string {
	if e.Reason == "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	} else {
		return fmt.Sprintf("%s (%s): %s", e.Code, e.Reason, e.Message)
	}
}

func (e *GenericError) WithMsg(msg string) *GenericError {
	return &GenericError{
		Code:     e.Code,
		Reason:   e.Reason,
		HttpCode: e.HttpCode,
		Message:  msg,
	}
}

func (e *GenericError) WithReason(reason string) *GenericError {
	return &GenericError{
		Code:     e.Code,
		Reason:   reason,
		HttpCode: e.HttpCode,
		Message:  e.Message,
	}
}
