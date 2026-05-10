package inmemory

func (s *SessionStorage) Save(jti string, session Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[jti] = session
}
