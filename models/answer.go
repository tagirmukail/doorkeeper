package models

import (
	"doorkeeper/utils"
	"fmt"
	"net/http"
)

type Answer struct {
	ID             string      `json:"id"`              // unique id
	Status         string      `json:"status"`          // http response status
	Headers        http.Header `json:"headers"`         // http response headers
	ResponseLength int64       `json:"response_length"` // http response length
}

func NewAnswer(resp *http.Response) (*Answer, error) {
	var id = utils.GenerateUUID()
	if id == "" {
		return nil, fmt.Errorf("GenerateUUID() returned empty id")
	}

	var answer = &Answer{
		ID:             id,
		Status:         resp.Status,
		Headers:        resp.Header,
		ResponseLength: resp.ContentLength,
	}

	return answer, nil
}
