package service

import (
	"fmt"
)

type StandardError struct {
	Code        uint
	ActualError error
	Line        string
	Message     string
}

func (s StandardError) Error() string {
	errStr := fmt.Sprintf("Code : %v Line:%v \n Error:%v \n Message: %v", s.Code, s.Line, s.ActualError, s.Message)
	return errStr
}

func getStandardError(err error, code uint, line, message string) *StandardError {
	return &StandardError{
		Code:        code,
		ActualError: err,
		Line:        line,
		Message:     message,
	}
}
