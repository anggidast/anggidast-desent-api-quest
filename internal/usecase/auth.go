package usecase

import (
	"strings"

	"desent-api-quest/internal/domain"
	"desent-api-quest/internal/token"
)

const (
	defaultUsername = "admin"
	defaultPassword = "secret"
)

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

	if username != defaultUsername || password != defaultPassword {
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
