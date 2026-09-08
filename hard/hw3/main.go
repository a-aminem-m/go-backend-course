package main

import (
	"net/http"
)

func main() {
	storage := NewMemoryStorage()
	http.HandleFunc("/task", createTaskHandler(storage))
	http.HandleFunc("/status/", getTaskStatusHandler(storage))
	http.HandleFunc("/result/", getTaskResultHandler(storage))
	http.HandleFunc("/register", registerHandler(storage))
	http.HandleFunc("/login", loginHandler(storage))
	http.HandleFunc("/commit", commitTaskHandler(storage))
	http.ListenAndServe(":8000", nil)
}
