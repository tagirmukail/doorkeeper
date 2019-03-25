package models

import (
	"fmt"
	"net/http"
	"net/url"
)

// Task - struct of represent task of request
type Task struct {
	Method  string `json:"method"`
	Address string `json:"address"`
}

func NewTask(method, adderss string) *Task {
	return &Task{
		Method:  method,
		Address: adderss,
	}
}

// validate task
func (t *Task) Validate() error {
	switch t.Method {
	case http.MethodGet:
		break
	case http.MethodPost:
		break
	case http.MethodPut:
		break
	case http.MethodDelete:
		break
	default:
		return fmt.Errorf("method not allowed, must be 'GET', 'POST','PUT','DELETE'")
	}

	_, err := url.ParseRequestURI(t.Address)
	if err != nil {
		return err
	}

	return nil
}
