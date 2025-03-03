package services

import (
	"context"
	"testwire/internal/dto/request"
)

type AuthenticationService interface {
	Register(request.CreateUserRequest) error
	Login(string, string) (string, string, error)
	Logout(context.Context, string, string) error
}
