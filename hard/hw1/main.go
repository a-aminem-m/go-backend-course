package main

import (
	"net/http"
)

func main() {
	storage := NewMemoryStorage()
	http.HandleFunc("/task", createTaskHandler(storage))
	http.HandleFunc("/status/", getTaskStatusHandler(storage))
	http.ListenAndServe(":8000", nil)
}
