package service

import (
	"context"
)

type CartService interface {
	Create(ctx context.Context, idOwner int64) (error, string)
	Get(ctx context.Context, idOwner int64) (error, string)
}
