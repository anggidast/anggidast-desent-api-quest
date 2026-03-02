package middleware

import (
	"net/http"
	"strings"

	"desent-api-quest/internal/domain"
	"desent-api-quest/internal/httpjson"
)

func RequireBearerToken(validate func(string) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := strings.TrimSpace(r.Header.Get("Authorization"))
			if header == "" {
				httpjson.WriteError(w, &domain.FieldError{Code: domain.ErrUnauthorized, Message: "missing bearer token"})
				return
			}

			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || strings.TrimSpace(token) == "" {
				httpjson.WriteError(w, &domain.FieldError{Code: domain.ErrUnauthorized, Message: "invalid authorization header"})
				return
			}

			if !validate(strings.TrimSpace(token)) {
				httpjson.WriteError(w, &domain.FieldError{Code: domain.ErrUnauthorized, Message: "invalid token"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
