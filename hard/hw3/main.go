package main

import (
	"log"
	"net/http"
)

func main() {
	db, err := NewPostgres()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rdb := NewRedis()
	defer rdb.Close()
	http.HandleFunc("/task", createTaskHandler(db, rdb))
	http.HandleFunc("/status/", getTaskStatusHandler(db, rdb))
	http.HandleFunc("/result/", getTaskResultHandler(db, rdb))
	http.HandleFunc("/register", registerHandler(db))
	http.HandleFunc("/login", loginHandler(db, rdb))
	http.HandleFunc("/commit", commitTaskHandler(db))
	http.ListenAndServe(":8000", nil)
}
