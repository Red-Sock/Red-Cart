package service

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type UserService interface {
	Start(ctx context.Context, user domain.User) (message string, err error)
}
