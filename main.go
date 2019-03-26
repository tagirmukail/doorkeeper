package main

import (
	"doorkeeper/handlers"
	"doorkeeper/worker"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"sync"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(string(debug.Stack()))
		}
	}()

	var addr string
	var workersCount int
	flag.StringVar(&addr, "addr", "0.0.0.0:8000", "Default: 0.0.0.0:8000")
	flag.IntVar(&workersCount, "workers", 2, "Default: 2")
	flag.Parse()

	log.Printf("<<<<<<<<Service started: %v>>>>>>>>", os.Args[1:])

	var wg = &sync.WaitGroup{}

	var taskWorker = worker.NewWorker(wg)
	go taskWorker.Run(workersCount)

	var router = mux.NewRouter()
	router.HandleFunc("/v1/fetch_task", handlers.FetchTask(taskWorker.TaskChan)).Methods("GET")

	log.Fatal(http.ListenAndServe(addr, router))
}
