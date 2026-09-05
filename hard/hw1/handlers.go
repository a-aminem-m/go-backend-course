package main

import (
	"fmt"
	"net/http"
)

func createTaskHandler(storage *MemoryStorage) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "task handler works")
	}
	return handler
}
