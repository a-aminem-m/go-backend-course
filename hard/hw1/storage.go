package main

import "sync"

type MemoryStorage struct {
	tasks      map[string]Task
	users      map[string]User
	sessions   map[string]Session
	tasksMu    sync.RWMutex
	usersMu    sync.RWMutex
	sessionsMu sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	storage := MemoryStorage{}
	storage.tasks = make(map[string]Task)
	storage.users = make(map[string]User)
	storage.sessions = make(map[string]Session)
	return &storage
}

func (s *MemoryStorage) CreateTask(task Task) {
	s.tasksMu.RLock()
	defer s.tasksMu.Unlock()
	s.tasks[task.ID] = task
}

func (s *MemoryStorage) GetTask(id string) (Task, bool) {
	s.tasksMu.RLock()
	defer s.tasksMu.RUnlock()
	task, ok := s.tasks[id]
	return task, ok
}

func (s *MemoryStorage) UpdateTask(task Task) {
	s.tasksMu.Lock()
	defer s.tasksMu.Unlock()
	s.tasks[task.ID] = task
}

func (s *MemoryStorage) CreateUser(user User) {
	s.usersMu.Lock()
	defer s.usersMu.Unlock()
	s.users[user.Login] = user
}

func (s *MemoryStorage) GetUser(login string) (User, bool) {
	s.usersMu.RLock()
	defer s.usersMu.RUnlock()
	user, ok := s.users[login]
	return user, ok
}

func (s *MemoryStorage) CreateSession(session Session) {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()
	s.sessions[session.SessionID] = session
}

func (s *MemoryStorage) GetSession(sessionID string) (Session, bool) {
	s.sessionsMu.RLock()
	defer s.sessionsMu.RUnlock()
	session, ok := s.sessions[sessionID]
	return session, ok
}
