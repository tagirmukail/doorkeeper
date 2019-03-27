package main

import (
	"doorkeeper/config"
	"doorkeeper/handlers"
	"doorkeeper/worker"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(string(debug.Stack()))
		}
	}()

	flag.Parse()

	cfg, err := config.NewConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("<<<<<<<<Service started: %+v>>>>>>>>", cfg)

	var wg = &sync.WaitGroup{}

	var taskWorker = worker.NewWorker(wg)
	go taskWorker.Run(cfg.Workers)

	var router = mux.NewRouter()
	router.HandleFunc("/v1/fetchtask", handlers.FetchTask(taskWorker.TaskChan)).Methods(http.MethodGet)
	router.HandleFunc("/v1/tasks/{page}", handlers.GetTasks(taskWorker, cfg.TaskCountOnPage)).Methods(http.MethodGet)
	router.HandleFunc("/v1/tasks/{id}", handlers.DeleteTask(taskWorker)).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(cfg.Address, router))
}
