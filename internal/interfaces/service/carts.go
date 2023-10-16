package service

import (
	"context"
)

type CartService interface {
	Create(ctx context.Context, idOwner int64) (string, error)
}
