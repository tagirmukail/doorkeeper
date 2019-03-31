package models

import (
	"doorkeeper/utils"
	"testing"
)

func TestTask_validate(t *testing.T) {
	type fields struct {
		ID      utils.UID
		Method  string
		Address string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Check valid task by method GET",
			fields: fields{
				ID:      utils.UID("test1"),
				Method:  "GET",
				Address: "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "Check valid task by method POST",
			fields: fields{
				ID:      utils.UID("test1"),
				Method:  "POST",
				Address: "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "Check valid task by method PUT",
			fields: fields{
				ID:      utils.UID("test1"),
				Method:  "PUT",
				Address: "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "Check valid task by method DELETE",
			fields: fields{
				ID:      utils.UID("test1"),
				Method:  "DELETE",
				Address: "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "Check invalid task by method",
			fields: fields{
				ID:      utils.UID("test1"),
				Method:  "Gt",
				Address: "http://example.com",
			},
			wantErr: true,
		},
		{
			name: "Check invalid task by address",
			fields: fields{
				ID:      utils.UID("test1"),
				Method:  "GET",
				Address: "example.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				ID:      tt.fields.ID,
				Method:  tt.fields.Method,
				Address: tt.fields.Address,
			}
			if err := task.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Task.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
