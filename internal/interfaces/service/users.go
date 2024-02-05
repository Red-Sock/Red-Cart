package service

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain/user"
)

type UserService interface {
	Start(ctx context.Context, user user.User) (message string, err error)
}
