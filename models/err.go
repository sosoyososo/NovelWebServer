package models

import (
	"fmt"
)

type ModelError struct {
	Code    int
	Message string
}

func (e ModelError) Error() string {
	return fmt.Sprintf("code :%c, message :%v", e.Message)
}

var (
	NullParameter = NewModelError(-1, "参数为空")
)

func NewModelError(code int, message string) ModelError {
	return ModelError{code, message}
}
