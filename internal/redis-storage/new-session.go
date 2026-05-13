package redisStorage

import "time"

type Session struct {
	UserID    string
	ExpiresAt time.Time
	Revoked   bool
}

func (r *Redis) NewSession(userID string) *Session {
	return &Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Minute * 15),
		Revoked:   false,
	}
}
