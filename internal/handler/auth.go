package handler

import (
	"net/http"

	"desent-api-quest/internal/httpjson"
	"desent-api-quest/internal/usecase"
)

type AuthHandler struct {
	auth *usecase.AuthUsecase
}

type tokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthHandler(auth *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) IssueToken(w http.ResponseWriter, r *http.Request) {
	var req tokenRequest
	if err := httpjson.DecodeJSON(r, &req); err != nil {
		httpjson.WriteError(w, err)
		return
	}

	token, err := h.auth.IssueToken(req.Username, req.Password)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}
