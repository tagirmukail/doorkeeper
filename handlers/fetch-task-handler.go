package handlers

import (
	"doorkeeper/models"
	"net/http"
)

func FetchTask(taskChan <-chan *models.Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			http.Error(w, "Request method allowed only GET", http.StatusForbidden)
			return
		}
	}
}
