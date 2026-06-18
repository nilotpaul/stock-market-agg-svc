package model

import "fmt"

type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("[%d]: %+v", e.Status, e.Message)
}
