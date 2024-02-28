package service

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type CartService interface {
	Create(ctx context.Context, idOwner int64) (string, error)

	UpdateMessageRef(ctx context.Context, item domain.Cart) error
}
