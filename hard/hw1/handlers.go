package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func createTaskHandler(storage *MemoryStorage) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !checkAuth(storage, r) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var task Task
		task.ID = uuid.New().String()
		task.Status = StatusInProgress
		task.Result = ""
		storage.CreateTask(task)
		go func(taskCopy Task) {
			time.Sleep(3 * time.Second)
			taskCopy.Status = StatusReady
			taskCopy.Result = "fake result"
			storage.UpdateTask(taskCopy)
		}(task)
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
		if !checkAuth(storage, r) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/status/")
		task, ok := storage.GetTask(id)
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
		if !checkAuth(storage, r) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/result/")
		task, ok := storage.GetTask(id)
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

func registerHandler(storage *MemoryStorage) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var request RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		var user User
		user.ID = uuid.New().String()
		user.Login = request.Username
		user.Password = request.Password
		storage.CreateUser(user)
		w.WriteHeader(http.StatusCreated)
	}
	return handler
}

func loginHandler(storage *MemoryStorage) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var request LoginRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		user, ok := storage.GetUser(request.Username)
		if !ok || request.Password != user.Password {
			http.Error(w, "user unauthorized", http.StatusUnauthorized)
			return
		}
		var session Session
		session.UserID = user.ID
		session.SessionID = uuid.New().String()
		storage.CreateSession(session)
		response := LoginResponse{
			Token: session.SessionID,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
	return handler
}

func checkAuth(storage *MemoryStorage, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	_, ok := storage.GetSession(token)
	return ok
}
