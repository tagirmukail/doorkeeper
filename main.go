package main

import (
	"flag"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"doorkeeper/config"
	"doorkeeper/handlers"
	"doorkeeper/worker"
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

	var tr = http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	var taskWorker = worker.NewWorker(wg, tr)
	go taskWorker.Run(cfg.Workers)

	var router = mux.NewRouter()
	router.HandleFunc("/v1/fetchtask", handlers.FetchTask(taskWorker)).Methods(http.MethodGet)
	router.HandleFunc("/v1/tasks/{page}", handlers.GetTasks(taskWorker, cfg.TaskCountOnPage)).Methods(http.MethodGet)
	router.HandleFunc("/v1/tasks/{id}", handlers.DeleteTask(taskWorker)).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(cfg.Address, router))
}
