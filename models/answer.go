package models

import "net/http"

type Answer struct {
	ID             string         `json:"id"`              // unique id
	Status         string         `json:"status"`          // http response status
	Headers        []*http.Header `json:"headers"`         // http response headers
	ResponseLength int            `json:"response_length"` // http response length
}
