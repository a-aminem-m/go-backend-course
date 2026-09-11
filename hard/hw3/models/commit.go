package models

type CommitRequest struct {
	TaskID string `json:"task_id"`
	Result string `json:"result"`
}
