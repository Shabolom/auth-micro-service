package inmemory

import (
	"context"
	"time"
)

func (s *SessionStorage) StartCleaner(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				s.DeleteExpired(time.Now())
			}
		}
	}()
}

func (s *SessionStorage) DeleteExpired(now time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	deleted := 0

	for jti, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			delete(s.sessions, jti)
			deleted++
		}
	}

	return deleted
}
