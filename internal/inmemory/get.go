package inmemory

import "fmt"

func (s *SessionStorage) Get(jti string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[jti]
	fmt.Println(session, ok)
	return session, ok
}
