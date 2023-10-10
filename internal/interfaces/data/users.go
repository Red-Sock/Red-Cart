package data

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain/user"
)

type Users interface {
	Upsert(ctx context.Context, u user.User) error
	Get(ctx context.Context, id int64) (user.User, error)
}
