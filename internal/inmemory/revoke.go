package inmemory

func (s *SessionStorage) Revoke(jti string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[jti]
	if !ok {
		return
	}

	session.Revoked = true
	s.sessions[jti] = session
}
