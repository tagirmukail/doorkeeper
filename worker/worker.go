package worker

import (
	"doorkeeper/models"
	"doorkeeper/utils"
	"log"
	"net/http"
	"sync"
)

type TaskCache []*models.Task // cache for tasks: map[UID]Task

// Worker of tasks
type Worker struct {
	taskChan   chan *models.Task // chanel for task
	wg         *sync.WaitGroup
	taskCache  TaskCache
	taskClient *http.Client
	*sync.RWMutex
}

func NewWorker(wg *sync.WaitGroup, tr http.Transport) *Worker {
	return &Worker{
		wg:         wg,
		RWMutex:    &sync.RWMutex{},
		taskChan:   make(chan *models.Task),
		taskCache:  []*models.Task{},
		taskClient: &http.Client{Transport: &tr},
	}
}

// Run - run gorutines of processing incoming tasks
func (w *Worker) Run(workers int) {
	for i := 0; i < workers; i++ {
		w.wg.Add(1)
		go w.work()
	}

	log.Printf("Worker run.")
	w.wg.Wait()
}

// work - processing incoming tasks
func (w *Worker) work() {
	w.wg.Done()
	for task := range w.taskChan {
		w.saveTask(task)
	}
}

// SendTask for further processing
func (w *Worker) SendTask(task *models.Task) {
	w.taskChan <- task
}

// saveTask - save task in cache
func (w *Worker) saveTask(task *models.Task) {
	w.Lock()
	w.taskCache = append(w.taskCache, task)
	w.Unlock()
}

// DoTask http request by task
func (w *Worker) DoTask(t *models.Task) (*http.Response, error) {
	req, err := http.NewRequest(t.Method, t.Address, nil)
	if err != nil {
		return nil, err
	}

	resp, err := w.taskClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetTaskCache() return all saving tasks
func (w *Worker) GetTaskCache() *TaskCache {
	w.RLock()
	var result = &w.taskCache
	w.RUnlock()
	return result
}

func (w *Worker) countAllTasks() int {
	w.RLock()
	var result = len(w.taskCache)
	w.RUnlock()
	return result
}

func (w *Worker) GetTasksPage(pageNumber, taskCountOnPage int) []*models.Task {
	var (
		result        []*models.Task
		countAllTasks = w.countAllTasks()
		start         = (pageNumber - 1) * taskCountOnPage
		stop          = start + taskCountOnPage
	)

	if pageNumber <= 0 {
		return result
	}

	if start > countAllTasks {
		return result
	}

	if stop > countAllTasks {
		stop = countAllTasks
	}

	w.RLock()
	result = w.taskCache[start:stop]
	w.RUnlock()

	return result
}

func (w *Worker) DeleteTask(id utils.UID) {
	var resultCache []*models.Task

	w.Lock()
	for _, task := range w.taskCache {
		if task.ID == id {
			continue
		}

		resultCache = append(resultCache, task)
	}
	w.taskCache = resultCache
	w.Unlock()
}

func (w *Worker) SetCache(cache TaskCache) {
	w.Lock()
	w.taskCache = cache
	w.Unlock()
}
