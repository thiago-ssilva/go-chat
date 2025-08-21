package service

import (
	"github.com/thiago-ssilva/zap/internal/ws"
)

type UserService struct {
	hub *ws.Hub
}

func NewUserService(hub *ws.Hub) *UserService {
	return &UserService{
		hub: hub,
	}
}

type UserValidationError struct {
	Code    string
	Message string
}

func (e UserValidationError) Error() string {
	return e.Message
}

const (
	ErrCodeUsernameTaken  = "USERNAME_TAKEN"
	ErrCodeUsernameLength = "USERNAME_LENGTH"
)

func (s *UserService) ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return UserValidationError{
			Code:    ErrCodeUsernameLength,
			Message: "username must be 3-20 characters",
		}
	}

	if s.hub.IsUsernameTaken(username) {
		return UserValidationError{
			Code:    ErrCodeUsernameTaken,
			Message: "username is already taken",
		}
	}

	return nil
}
