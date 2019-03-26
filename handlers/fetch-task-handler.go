package handlers

import (
	"doorkeeper/models"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func FetchTask(taskChan chan<- *models.Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			http.Error(w, "Request method allowed only GET", http.StatusForbidden)
			return
		}

		taskStr := r.URL.Query().Get("task")
		if taskStr == "" {
			http.Error(w, "Task is empty", http.StatusBadRequest)
			return
		}

		task := &models.Task{}
		err := json.Unmarshal([]byte(taskStr), task)
		if err != nil {
			http.Error(w, "Task is not json", http.StatusBadRequest)
			return
		}

		err = task.Validate()
		if err != nil {
			log.Printf("FetchTask() validation error: %v", err)
			http.Error(w, "Task not valid", http.StatusBadRequest)
			return
		}

		taskChan <- task

		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		}

		httpClient := http.Client{Transport: tr}

		req, err := http.NewRequest(task.Method, task.Address, nil)
		if err != nil {
			log.Printf("FetchTask() NewRequest() error: %v", err)
			http.Error(w, "Task request error", http.StatusInternalServerError)
			return
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Printf("FetchTask() httpClient.Do() error: %v", err)
			http.Error(w, "Task request error", http.StatusInternalServerError)
			return
		}

		answer, err := models.NewAnswer(resp)
		if err != nil {
			log.Printf("FetchTask() NewAnswer error: %v", err)
			http.Error(w, "Answer error", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(answer)
		if err != nil {
			log.Printf("FetchTask() Encode Answer error: %v", err)
			http.Error(w, "Answer error", http.StatusInternalServerError)
			return
		}

		return
	}
}
