package models

import (
	"net/http"

	"doorkeeper/utils"
)

type Answer struct {
	ID             utils.UID   `json:"id"`              // unique id
	Status         string      `json:"status"`          // http response status
	Headers        http.Header `json:"headers"`         // http response headers
	ResponseLength int64       `json:"response_length"` // http response length
}

func NewAnswer(resp *http.Response, id utils.UID) *Answer {
	return &Answer{
		ID:             id,
		Status:         resp.Status,
		Headers:        resp.Header,
		ResponseLength: resp.ContentLength,
	}
}
