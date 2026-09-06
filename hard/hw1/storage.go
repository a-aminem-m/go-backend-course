package main

import "sync"

type MemoryStorage struct {
	tasks    map[string]Task
	users    map[string]User
	sessions map[string]Session
	mu       sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	storage := MemoryStorage{}
	storage.tasks = make(map[string]Task)
	storage.users = make(map[string]User)
	storage.sessions = make(map[string]Session)
	return &storage
}

func (s *MemoryStorage) CreateTask(task Task) {
	s.tasks[task.ID] = task
}

func (s *MemoryStorage) GetTask(id string) (Task, bool) {
	task, ok := s.tasks[id]
	return task, ok
}

func (s *MemoryStorage) UpdateTask(task Task) {
	s.tasks[task.ID] = task
}

func (s *MemoryStorage) CreateUser(user User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.Login] = user
}

func (s *MemoryStorage) GetUser(login string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[login]
	return user, ok
}

func (s *MemoryStorage) CreateSession(session Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.SessionID] = session
}

func (s *MemoryStorage) GetSession(sessionID string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[sessionID]
	return session, ok
}
