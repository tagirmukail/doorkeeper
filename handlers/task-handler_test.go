package handlers

import (
	"doorkeeper/models"
	"doorkeeper/utils"
	"doorkeeper/worker"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

func TestGetTasks(t *testing.T) {
	type args struct {
		worker          *worker.Worker
		taskCountOnPage int
	}
	tests := []struct {
		name       string
		args       args
		taskCache  worker.TaskCache
		wantResp   []*models.Task
		wantStatus int
	}{
		{
			name: "Get tasks by page number",
			args: args{
				worker: worker.NewWorker(
					&sync.WaitGroup{},
					http.Transport{
						IdleConnTimeout:    4,
						MaxIdleConns:       10,
						DisableCompression: true,
					},
				),
				taskCountOnPage: 3,
			},
			taskCache: []*models.Task{
				{
					Address: "http://test.com",
					Method:  "GET",
					ID:      utils.UID("test1"),
				},
			},
			wantResp: []*models.Task{
				{
					Address: "http://test.com",
					Method:  "GET",
					ID:      utils.UID("test1"),
				},
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.worker.SetCache(tt.taskCache)

			req, err := http.NewRequest("GET", "/v1/tasks/1", nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{"page": "1"})

			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(GetTasks(tt.args.worker, tt.args.taskCountOnPage))
			handler.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tt.wantStatus {
				t.Errorf("GetTasks() handler wrong status code: got %v want %v", status, tt.wantStatus)
			}

			var gotTasks []*models.Task
			body, _ := ioutil.ReadAll(recorder.Body)
			if err != nil {
				t.Error(err)
			}
			err = json.Unmarshal(body, &gotTasks)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(tt.wantResp, gotTasks) {
				t.Errorf("GetTasks() handler wrong response: got %+v want %+v", tt.wantResp, gotTasks)
			}
		})
	}
}
