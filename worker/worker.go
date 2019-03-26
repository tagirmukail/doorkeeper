package worker

import (
	"doorkeeper/models"
	"sync"
)

type Worker struct {
	TaskChan chan *models.Task // chanel for task
	wg       *sync.WaitGroup

	*sync.RWMutex
	taskCache map[string]string // cache for tasks: map[address]method
}

func NewWorker(wg *sync.WaitGroup) *Worker {
	return &Worker{
		wg:        wg,
		RWMutex:   &sync.RWMutex{},
		TaskChan:  make(chan *models.Task),
		taskCache: make(map[string]string),
	}
}

func (w *Worker) Run(workers int) {
	for i := 0; i < workers; i++ {
		w.wg.Add(1)
		go w.work()
	}

	w.wg.Wait()
}

func (w *Worker) work() {
	w.wg.Done()
	for task := range w.TaskChan {
		w.saveTask(task)
	}
}

func (w *Worker) saveTask(task *models.Task) {
	w.Lock()
	w.taskCache[task.Address] = task.Method
	w.Unlock()
}
