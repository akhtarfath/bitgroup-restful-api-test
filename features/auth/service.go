package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/akhtarfath/config"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type TokenService struct {
	tokens map[string]time.Time
	mu     sync.Mutex
}

func NewTokenService() *TokenService {
	return &TokenService{
		tokens: make(map[string]time.Time),
	}
}

func (s *TokenService) Authenticate(username, password string) (string, time.Time, error) {
	if username != config.Username() || password != config.Password() {
		return "", time.Time{}, ErrInvalidCredentials
	}

	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", time.Time{}, err
	}

	token := hex.EncodeToString(buf)
	expiresAt := time.Now().Add(config.TokenTTL())

	s.mu.Lock()
	s.tokens[token] = expiresAt
	s.mu.Unlock()

	return token, expiresAt, nil
}

func (s *TokenService) IsValid(token string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	expiresAt, ok := s.tokens[token]
	if !ok {
		return false
	}
	if time.Now().After(expiresAt) {
		delete(s.tokens, token)
		return false
	}

	return true
}
