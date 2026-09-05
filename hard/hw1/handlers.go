package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type PostTaskResponse struct {
	TaskID string `json:"task_id"`
}

type GetTaskStatusResponse struct {
	Status string `json:"status"`
}

type GetTaskResultResponse struct {
	Result string `json:"result"`
}

func createTaskHandler(storage *MemoryStorage) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var task Task
		task.ID = uuid.New().String()
		task.Status = StatusInProgress
		task.Result = ""
		storage.Create(task)
		go func() {
			time.Sleep(3 * time.Second)
			task.Status = StatusReady
			task.Result = "fake result"
			storage.Update(task)
		}()
		response := PostTaskResponse{
			TaskID: task.ID,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
	return handler
}

func getTaskStatusHandler(storage *MemoryStorage) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/status/")
		task, ok := storage.Get(id)
		if !ok {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		response := GetTaskStatusResponse{
			Status: task.Status,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
	return handler
}

func getTaskResultHandler(storage *MemoryStorage) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/result/")
		task, ok := storage.Get(id)
		if !ok {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		response := GetTaskResultResponse{
			Result: task.Result,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
	return handler
}
