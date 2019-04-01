package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"doorkeeper/models"
	"doorkeeper/utils"
	"doorkeeper/worker"
)

func FetchTask(taskWorker *worker.Worker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		taskStr := r.URL.Query().Get("task")
		if taskStr == "" {
			http.Error(w, "Task is empty", http.StatusBadRequest)
			return
		}

		task, err := models.NewTask([]byte(taskStr))
		if err != nil {
			log.Printf("models.NewTask() error: %v", err)
			http.Error(w, "Bad task", http.StatusBadRequest)
			return
		}

		taskWorker.SendTask(task)

		resp, err := taskWorker.DoTask(task)
		if err != nil {
			log.Printf("FetchTask() error: %v", err)
			http.Error(w, "Task request error", http.StatusBadRequest)
			return
		}

		answer := models.NewAnswer(resp, task.ID)
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
