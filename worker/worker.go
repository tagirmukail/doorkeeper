package worker

import (
	"log"
	"sync"

	"doorkeeper/models"
	"doorkeeper/utils"
)

type TaskCache map[utils.UID]*models.Task // cache for tasks: map[UID]Task

// Worker of tasks
type Worker struct {
	TaskChan  chan *models.Task // chanel for task
	wg        *sync.WaitGroup
	taskCache TaskCache
	*sync.RWMutex
}

func NewWorker(wg *sync.WaitGroup) *Worker {
	return &Worker{
		wg:        wg,
		RWMutex:   &sync.RWMutex{},
		TaskChan:  make(chan *models.Task),
		taskCache: make(map[utils.UID]*models.Task),
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
	for task := range w.TaskChan {
		w.saveTask(task)
	}
}

// saveTask - save task in cache
func (w *Worker) saveTask(task *models.Task) {
	w.Lock()
	w.taskCache[task.ID] = task
	w.Unlock()
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
		count         int
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
	for _, task := range w.taskCache {
		if count < start || count >= stop {
			count++
			continue
		}

		count++
		result = append(result, task)
	}
	w.RUnlock()

	return result
}

func (w *Worker) DeleteTask(id utils.UID) {
	w.Lock()
	delete(w.taskCache, id)
	w.Unlock()
}
