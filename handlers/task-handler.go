package handlers

import (
	"doorkeeper/models"
	"doorkeeper/utils"
	"doorkeeper/worker"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func FetchTask(taskChan chan<- *models.Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

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

		uid, err := utils.GenerateUID()
		if err != nil {
			log.Printf("FetchTask() GenerateUID() error: %v", err)
			http.Error(w, "Task not valid", http.StatusInternalServerError)
			return
		}

		task.ID = uid

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

		answer := models.NewAnswer(resp, uid)
		err = json.NewEncoder(w).Encode(answer)
		if err != nil {
			log.Printf("FetchTask() Encode Answer error: %v", err)
			http.Error(w, "Answer error", http.StatusInternalServerError)
			return
		}

		return
	}
}

func GetTasks(worker *worker.Worker, taskCountOnPage int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var pageNumberStr = mux.Vars(r)["page"]
		pageNumber, err := strconv.Atoi(pageNumberStr)
		if err != nil {
			http.Error(w, "Page number must be greater than 0", http.StatusForbidden)
			return
		}

		var tasks = worker.GetTasksPage(pageNumber, taskCountOnPage)
		if len(tasks) == 0 {
			http.Error(w, "Not found tasks", http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(tasks)
		if err != nil {
			log.Printf("GetTasks() Encode Tasks error: %v", err)
			http.Error(w, "Get Tasks error", http.StatusInternalServerError)
			return
		}

		return
	}
}

func DeleteTask(worker *worker.Worker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var uid = mux.Vars(r)["id"]

		worker.DeleteTask(utils.UID(uid))

		w.WriteHeader(http.StatusOK)
		return
	}
}
