package main

const (
	StatusInProgress = "in_progress"
	StatusReady      = "ready"
)

type Task struct {
	ID     string
	Status string
	Result string
}
