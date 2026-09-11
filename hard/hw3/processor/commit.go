package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hard-hw1/models"
	"net/http"
	"os"
)

func sendResult(taskID string, result string) error {
	request := models.CommitRequest{
		TaskID: taskID,
		Result: result,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8000"
	}
	response, err := http.Post(
		serverURL+"/commit",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("commit failed with status: %s", response.Status)
	}

	return nil
}
