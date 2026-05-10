package inmemory

func (s *SessionStorage) Get(jti string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[jti]
	return session, ok
}
