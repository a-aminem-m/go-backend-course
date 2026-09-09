package main

import (
	"database/sql"
	"encoding/json"
	"hard-hw1/models"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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

func createTaskHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !checkAuth(rdb, r) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var task Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		task.ID = uuid.New().String()
		task.Status = StatusInProgress
		task.Result = ""
		err = CreateTask(db, task)
		if err != nil {
			http.Error(w, "failed to create task", http.StatusInternalServerError)
			return
		}
		message := models.CodeTaskMessage{
			TaskID:     task.ID,
			Translator: task.Translator,
			Code:       task.Code,
		}
		err = publishTask(message)
		if err != nil {
			http.Error(w, "failed to publish task", http.StatusInternalServerError)
			return
		}
		response := PostTaskResponse{
			TaskID: task.ID,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
	return handler
}

func getTaskStatusHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !checkAuth(rdb, r) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/status/")
		task, ok := GetTask(db, id)
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

func getTaskResultHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !checkAuth(rdb, r) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/result/")
		task, ok := GetTask(db, id)
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

func registerHandler(db *sql.DB) http.HandlerFunc {
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
		err = CreateUser(db, user)
		if err != nil {
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
	return handler
}

func loginHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
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
		user, ok := GetUser(db, request.Username)
		if !ok || request.Password != user.Password {
			http.Error(w, "user unauthorized", http.StatusUnauthorized)
			return
		}
		var session Session
		session.UserID = user.ID
		session.SessionID = uuid.New().String()
		err = CreateSession(rdb, session)
		if err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}
		response := LoginResponse{
			Token: session.SessionID,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
	return handler
}

func checkAuth(rdb *redis.Client, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	_, ok := GetSession(rdb, token)
	return ok
}

func commitTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request models.CommitRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		task, ok := GetTask(db, request.TaskID)
		if !ok {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}

		task.Result = request.Result
		task.Status = StatusReady

		err = UpdateTask(db, task)
		if err != nil {
			http.Error(w, "failed to update task", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
