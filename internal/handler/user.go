package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiago-ssilva/zap/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) ValidateUsername(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	username := q.Get("username")

	response := struct {
		Error   string `json:"error,omitempty"`
		ErrCode string `json:"code,omitempty"`
		Valid   bool   `json:"valid"`
	}{
		ErrCode: "",
		Error:   "",
		Valid:   true,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := h.userService.ValidateUsername(username); err != nil {
		if validationErr, ok := err.(service.UserValidationError); ok {
			w.WriteHeader(http.StatusBadRequest)
			response.Valid = false
			response.Error = validationErr.Message
			response.ErrCode = validationErr.Code
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}

	_ = json.NewEncoder(w).Encode(response)
}
