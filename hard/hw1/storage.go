package main

import "sync"

type MemoryStorage struct {
	tasks map[string]Task
	mu    sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	storage := MemoryStorage{}
	storage.tasks = make(map[string]Task)
	return &storage
}

func (s *MemoryStorage) Create(task Task) {
	s.tasks[task.ID] = task
}

func (s *MemoryStorage) Get(id string) (Task, bool) {
	task, ok := s.tasks[id]
	return task, ok
}

func (s *MemoryStorage) Update(task Task) {
	s.tasks[task.ID] = task
}
