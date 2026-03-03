package usecase

import (
	"strings"

	"desent-api-quest/internal/domain"
	"desent-api-quest/internal/token"
)

var validCredentialPairs = map[string]map[string]struct{}{
	"admin": {
		"secret":   {},
		"password": {},
	},
	"user": {
		"password": {},
	},
}

type AuthUsecase struct {
	tokens *token.Service
}

func NewAuthUsecase(tokens *token.Service) *AuthUsecase {
	return &AuthUsecase{tokens: tokens}
}

func (u *AuthUsecase) IssueToken(username, password string) (string, error) {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return "", &domain.FieldError{Code: domain.ErrBadRequest, Message: "username and password are required"}
	}

	passwords, ok := validCredentialPairs[username]
	if !ok {
		return "", &domain.FieldError{Code: domain.ErrUnauthorized, Message: "invalid credentials"}
	}
	if _, ok := passwords[password]; !ok {
		return "", &domain.FieldError{Code: domain.ErrUnauthorized, Message: "invalid credentials"}
	}

	token, err := u.tokens.Issue()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *AuthUsecase) ValidateToken(value string) bool {
	return u.tokens.Validate(value)
}
