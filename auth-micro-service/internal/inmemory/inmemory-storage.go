package inmemory

import (
	"sync"
	"time"
)

type Session struct {
	UserID    string
	ExpiresAt time.Time
	Revoked   bool
}

type SessionStorage struct {
	mu       sync.RWMutex
	sessions map[string]Session // key = jti
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		sessions: make(map[string]Session),
	}
}

func (s *SessionStorage) NewSession(userID string) Session {

	return Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Minute * 15),
		Revoked:   false,
	}
}
