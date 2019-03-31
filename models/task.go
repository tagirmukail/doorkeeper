package models

import (
	"doorkeeper/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Task - struct of represent task of request
type Task struct {
	ID      utils.UID `json:"id"`
	Method  string    `json:"method"`
	Address string    `json:"address"`
}

func NewTask(taskByte []byte) (*Task, error) {
	var task = &Task{}

	err := json.Unmarshal([]byte(taskByte), task)
	if err != nil {
		return nil, err
	}

	err = task.validate()
	if err != nil {
		return nil, err
	}

	uid, err := utils.GenerateUID()
	if err != nil {
		return nil, err
	}

	task.ID = uid
	return task, nil
}

// validate task
func (t *Task) validate() error {
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
