package token

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type Service struct {
	mu     sync.RWMutex
	tokens map[string]struct{}
}

func NewService() *Service {
	return &Service{
		tokens: make(map[string]struct{}),
	}
}

func (s *Service) Issue() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	token := hex.EncodeToString(raw)

	s.mu.Lock()
	s.tokens[token] = struct{}{}
	s.mu.Unlock()

	return token, nil
}

func (s *Service) Validate(token string) bool {
	s.mu.RLock()
	_, ok := s.tokens[token]
	s.mu.RUnlock()
	return ok
}
