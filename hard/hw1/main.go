package main

import (
	"net/http"
)

func main() {
	storage := NewMemoryStorage()
	_ = storage
	http.ListenAndServe(":8000", nil)
}
