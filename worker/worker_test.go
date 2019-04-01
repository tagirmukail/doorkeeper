package worker

import (
	"doorkeeper/models"
	"doorkeeper/utils"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

func TestWorker_saveTask(t *testing.T) {
	type fields struct {
		taskChan  chan *models.Task
		wg        *sync.WaitGroup
		taskCache TaskCache
		RWMutex   *sync.RWMutex
	}
	type args struct {
		task *models.Task
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Check saving task 1",
			fields: fields{
				taskChan:  make(chan *models.Task),
				wg:        &sync.WaitGroup{},
				RWMutex:   &sync.RWMutex{},
				taskCache: TaskCache{},
			},
			args: args{
				task: &models.Task{
					ID:      utils.UID("testuid"),
					Address: "http://example.com",
					Method:  "GET",
				},
			},
		},
		{
			name: "Check saving task 2",
			fields: fields{
				taskChan:  make(chan *models.Task),
				wg:        &sync.WaitGroup{},
				RWMutex:   &sync.RWMutex{},
				taskCache: TaskCache{},
			},
			args: args{
				task: &models.Task{
					ID:      utils.UID("testuid"),
					Address: "http://example.com",
					Method:  "GET",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				taskChan:  tt.fields.taskChan,
				wg:        tt.fields.wg,
				taskCache: tt.fields.taskCache,
				RWMutex:   tt.fields.RWMutex,
			}
			w.saveTask(tt.args.task)

			w.RLock()
			var savedTask *models.Task
			for _, task := range w.taskCache {
				if task.ID == tt.args.task.ID {
					savedTask = task
				}
			}
			w.RUnlock()
			if savedTask != tt.args.task {
				t.Errorf("Saved task %+v not equal want task %+v", savedTask, tt.args.task)
			}
		})
	}
}

func TestWorker_GetTaskCache(t *testing.T) {
	type fields struct {
		taskChan  chan *models.Task
		wg        *sync.WaitGroup
		taskCache TaskCache
		RWMutex   *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   TaskCache
	}{
		{
			name: "Check task cache",
			fields: fields{
				taskChan: make(chan *models.Task),
				wg:       &sync.WaitGroup{},
				RWMutex:  &sync.RWMutex{},
				taskCache: TaskCache{
					{
						ID:      utils.UID("test1"),
						Address: "http://example.com",
						Method:  "GET",
					},
					{
						ID:      utils.UID("test2"),
						Address: "http://example2.com",
						Method:  "POST",
					},
				},
			},
			want: TaskCache([]*models.Task{
				{
					ID:      utils.UID("test1"),
					Address: "http://example.com",
					Method:  "GET",
				},
				{
					ID:      utils.UID("test2"),
					Address: "http://example2.com",
					Method:  "POST",
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				taskChan:  tt.fields.taskChan,
				wg:        tt.fields.wg,
				taskCache: tt.fields.taskCache,
				RWMutex:   tt.fields.RWMutex,
			}
			if got := w.GetTaskCache(); !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("Worker.GetTaskCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorker_countAllTasks(t *testing.T) {
	type fields struct {
		taskChan  chan *models.Task
		wg        *sync.WaitGroup
		taskCache TaskCache
		RWMutex   *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Check count tasks in tasks cache",
			fields: fields{
				taskChan: make(chan *models.Task),
				wg:       &sync.WaitGroup{},
				RWMutex:  &sync.RWMutex{},
				taskCache: TaskCache{
					{
						ID:      utils.UID("test1"),
						Address: "http://example.com",
						Method:  "GET",
					},
					{
						ID:      utils.UID("test2"),
						Address: "http://example2.com",
						Method:  "POST",
					},
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				taskChan:  tt.fields.taskChan,
				wg:        tt.fields.wg,
				taskCache: tt.fields.taskCache,
				RWMutex:   tt.fields.RWMutex,
			}
			if got := w.countAllTasks(); got != tt.want {
				t.Errorf("Worker.countAllTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorker_GetTasksPage(t *testing.T) {
	type fields struct {
		taskChan  chan *models.Task
		wg        *sync.WaitGroup
		taskCache TaskCache
		RWMutex   *sync.RWMutex
	}
	type args struct {
		pageNumber      int
		taskCountOnPage int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*models.Task
	}{
		{
			name: "Check task cache for first page",
			fields: fields{
				taskChan: make(chan *models.Task),
				wg:       &sync.WaitGroup{},
				RWMutex:  &sync.RWMutex{},
				taskCache: TaskCache{
					{
						ID:      utils.UID("test1"),
						Address: "http://example1.com",
						Method:  "GET",
					},
					{
						ID:      utils.UID("test2"),
						Address: "http://example2.com",
						Method:  "POST",
					},
					{
						ID:      utils.UID("test3"),
						Address: "http://example3.com",
						Method:  "PUT",
					},
					{
						ID:      utils.UID("test4"),
						Address: "http://example4.com",
						Method:  "POST",
					},
					{
						ID:      utils.UID("test5"),
						Address: "http://example5.com",
						Method:  "GET",
					},
					{
						ID:      utils.UID("test6"),
						Address: "http://example6.com",
						Method:  "POST",
					},
					{
						ID:      utils.UID("test7"),
						Address: "http://example7.com",
						Method:  "GET",
					},
					{
						ID:      utils.UID("test8"),
						Address: "http://example8.com",
						Method:  "DELETE",
					},
					{
						ID:      utils.UID("test9"),
						Address: "http://example9.com",
						Method:  "POST",
					},
					{
						ID:      utils.UID("test10"),
						Address: "http://example10.com",
						Method:  "PUT",
					},
					{
						ID:      utils.UID("test11"),
						Address: "http://example11.com",
						Method:  "DELETE",
					},
					{
						ID:      utils.UID("test12"),
						Address: "http://example12.com",
						Method:  "GET",
					},
				},
			},
			args: args{
				pageNumber:      1,
				taskCountOnPage: 5,
			},
			want: []*models.Task{
				{
					ID:      utils.UID("test1"),
					Address: "http://example1.com",
					Method:  "GET",
				},
				{
					ID:      utils.UID("test2"),
					Address: "http://example2.com",
					Method:  "POST",
				},
				{
					ID:      utils.UID("test3"),
					Address: "http://example3.com",
					Method:  "PUT",
				},
				{
					ID:      utils.UID("test4"),
					Address: "http://example4.com",
					Method:  "POST",
				},
				{
					ID:      utils.UID("test5"),
					Address: "http://example5.com",
					Method:  "GET",
				},
			},
		},
		{
			name: "Check task cache for second page",
			fields: fields{
				taskChan: make(chan *models.Task),
				wg:       &sync.WaitGroup{},
				RWMutex:  &sync.RWMutex{},
				taskCache: TaskCache{
					{
						ID:      utils.UID("test1"),
						Address: "http://example1.com",
						Method:  "GET",
					},
					{
						ID:      utils.UID("test2"),
						Address: "http://example2.com",
						Method:  "POST",
					},
					{
						ID:      utils.UID("test3"),
						Address: "http://example3.com",
						Method:  "PUT",
					},
					{
						ID:      utils.UID("test4"),
						Address: "http://example4.com",
						Method:  "POST",
					},
					{
						ID:      utils.UID("test5"),
						Address: "http://example5.com",
						Method:  "GET",
					},
					{
						ID:      utils.UID("test6"),
						Address: "http://example6.com",
						Method:  "POST",
					},
					{
						ID:      utils.UID("test7"),
						Address: "http://example7.com",
						Method:  "GET",
					},
					{
						ID:      utils.UID("test8"),
						Address: "http://example8.com",
						Method:  "DELETE",
					},
					{
						ID:      utils.UID("test9"),
						Address: "http://example9.com",
						Method:  "POST",
					},
					{
						ID:      utils.UID("test10"),
						Address: "http://example10.com",
						Method:  "PUT",
					},
					{
						ID:      utils.UID("test11"),
						Address: "http://example11.com",
						Method:  "DELETE",
					},
					{
						ID:      utils.UID("test12"),
						Address: "http://example12.com",
						Method:  "GET",
					},
				},
			},
			args: args{
				pageNumber:      2,
				taskCountOnPage: 5,
			},
			want: []*models.Task{
				{
					ID:      utils.UID("test6"),
					Address: "http://example6.com",
					Method:  "POST",
				},
				{
					ID:      utils.UID("test7"),
					Address: "http://example7.com",
					Method:  "GET",
				},
				{
					ID:      utils.UID("test8"),
					Address: "http://example8.com",
					Method:  "DELETE",
				},
				{
					ID:      utils.UID("test9"),
					Address: "http://example9.com",
					Method:  "POST",
				},
				{
					ID:      utils.UID("test10"),
					Address: "http://example10.com",
					Method:  "PUT",
				},
			},
		},
		{
			name: "Check task cache for end page",
			fields: fields{
				taskChan: make(chan *models.Task),
				wg:       &sync.WaitGroup{},
				RWMutex:  &sync.RWMutex{},
				taskCache: TaskCache{
					{
						ID:      utils.UID("test1"),
						Address: "http://example1.com",
						Method:  "GET",
					},
					{
						ID:      utils.UID("test2"),
						Address: "http://example2.com",
						Method:  "POST",
					},
					{
						ID:      utils.UID("test3"),
						Address: "http://example3.com",
						Method:  "PUT",
					},
					{
						ID:      utils.UID("test4"),
						Address: "http://example4.com",
						Method:  "POST",
					},
					{
						ID:      utils.UID("test5"),
						Address: "http://example5.com",
						Method:  "GET",
					},
					{
						ID:      utils.UID("test6"),
						Address: "http://example6.com",
						Method:  "POST",
					},
					{
						ID:      utils.UID("test7"),
						Address: "http://example7.com",
						Method:  "GET",
					},
					{
						ID:      utils.UID("test8"),
						Address: "http://example8.com",
						Method:  "DELETE",
					},
					{
						ID:      utils.UID("test9"),
						Address: "http://example9.com",
						Method:  "POST",
					},
					{
						ID:      utils.UID("test10"),
						Address: "http://example10.com",
						Method:  "PUT",
					},
					{
						ID:      utils.UID("test11"),
						Address: "http://example11.com",
						Method:  "DELETE",
					},
					{
						ID:      utils.UID("test12"),
						Address: "http://example12.com",
						Method:  "GET",
					},
				},
			},
			args: args{
				pageNumber:      3,
				taskCountOnPage: 5,
			},
			want: []*models.Task{
				{
					ID:      utils.UID("test11"),
					Address: "http://example11.com",
					Method:  "DELETE",
				},
				{
					ID:      utils.UID("test12"),
					Address: "http://example12.com",
					Method:  "GET",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				taskChan:  tt.fields.taskChan,
				wg:        tt.fields.wg,
				taskCache: tt.fields.taskCache,
				RWMutex:   tt.fields.RWMutex,
			}

			if got := w.GetTasksPage(tt.args.pageNumber, tt.args.taskCountOnPage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Worker.GetTasksPage() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestWorker_DeleteTask(t *testing.T) {
	type fields struct {
		taskChan  chan *models.Task
		wg        *sync.WaitGroup
		taskCache TaskCache
		RWMutex   *sync.RWMutex
	}
	type args struct {
		id utils.UID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "1: Check delete task from tasks cache",
			fields: fields{
				taskChan: make(chan *models.Task),
				wg:       &sync.WaitGroup{},
				RWMutex:  &sync.RWMutex{},
				taskCache: TaskCache{
					{
						ID:      utils.UID("test1"),
						Address: "http://example.com",
						Method:  "GET",
					},
					{
						ID:      utils.UID("test2"),
						Address: "http://example2.com",
						Method:  "POST",
					},
				},
			},
			args: args{
				id: utils.UID("test1"),
			},
		},
		{
			name: "2: Check delete task from tasks cache",
			fields: fields{
				taskChan: make(chan *models.Task),
				wg:       &sync.WaitGroup{},
				RWMutex:  &sync.RWMutex{},
				taskCache: TaskCache{
					{
						ID:      utils.UID("test1"),
						Address: "http://example.com",
						Method:  "GET",
					},
					{
						ID:      utils.UID("test2"),
						Address: "http://example2.com",
						Method:  "POST",
					},
				},
			},
			args: args{
				id: utils.UID("test2"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				taskChan:  tt.fields.taskChan,
				wg:        tt.fields.wg,
				taskCache: tt.fields.taskCache,
				RWMutex:   tt.fields.RWMutex,
			}
			w.DeleteTask(tt.args.id)

			for _, task := range w.taskCache {
				if task.ID == tt.args.id {
					t.Errorf("Task %v must be deleted", tt.args.id)
				}
			}
		})
	}
}

func TestWorker_DoTask(t *testing.T) {
	type fields struct {
		taskChan   chan *models.Task
		wg         *sync.WaitGroup
		taskCache  TaskCache
		taskClient *http.Client
		RWMutex    *sync.RWMutex
	}
	type args struct {
		t *models.Task
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantHeaderContentType string
		WantStatusCode        int
		wantErr               bool
	}{
		{
			name: "Check task http client",
			fields: fields{
				taskChan:  make(chan *models.Task),
				wg:        &sync.WaitGroup{},
				taskCache: TaskCache{},
				RWMutex:   &sync.RWMutex{},
			},
			args: args{
				t: &models.Task{
					ID:     utils.UID("test1"),
					Method: http.MethodGet,
				},
			},
			wantHeaderContentType: "application/json",
			WantStatusCode:        http.StatusOK,
			wantErr:               false,
		},
		{
			name: "2: Check task http client",
			fields: fields{
				taskChan:  make(chan *models.Task),
				wg:        &sync.WaitGroup{},
				taskCache: TaskCache{},
				RWMutex:   &sync.RWMutex{},
			},
			args: args{
				t: &models.Task{
					ID:     utils.UID("test1"),
					Method: http.MethodGet,
				},
			},
			wantHeaderContentType: "text",
			WantStatusCode:        http.StatusNotFound,
			wantErr:               false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", tt.wantHeaderContentType)
					w.WriteHeader(tt.WantStatusCode)
				}),
			)
			defer server.Close()

			tt.args.t.Address = server.URL

			w := &Worker{
				taskChan:   tt.fields.taskChan,
				wg:         tt.fields.wg,
				taskCache:  tt.fields.taskCache,
				taskClient: server.Client(),
				RWMutex:    tt.fields.RWMutex,
			}
			got, err := w.DoTask(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Worker.DoTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.StatusCode, tt.WantStatusCode) {
				t.Errorf("Worker.DoTask() = %v, wantStatus %v", got.StatusCode, tt.WantStatusCode)
			}

			if !reflect.DeepEqual(got.Header.Get("Content-Type"), tt.wantHeaderContentType) {
				t.Errorf("Worker.DoTask() = %v, wantContentType %v", got.Header.Get("Content-Type"), tt.wantHeaderContentType)
			}

			if err := got.Body.Close(); err != nil {
				t.Errorf("Worker.DoTask() got.Body.Close() error: %v", err)
			}
		})
	}
}
