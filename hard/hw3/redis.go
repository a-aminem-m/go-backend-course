package main

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func CreateSession(rdb *redis.Client, session Session) error {
	ctx := context.Background()

	err := rdb.Set(
		ctx,
		"session:"+session.SessionID,
		session.UserID,
		0,
	).Err()

	if err != nil {
		return err
	}

	return nil
}

func GetSession(rdb *redis.Client, sessionID string) (Session, bool) {
	ctx := context.Background()

	userID, err := rdb.Get(
		ctx,
		"session:"+sessionID,
	).Result()

	if err == redis.Nil {
		return Session{}, false
	}

	if err != nil {
		return Session{}, false
	}

	return Session{
		UserID:    userID,
		SessionID: sessionID,
	}, true
}
