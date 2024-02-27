package service

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type CartService interface {
	Create(ctx context.Context, idOwner int64) (string, error)
	GetByOwnerId(ctx context.Context, idOwner int64) (domain.Cart, error)
	GetByCartId(ctx context.Context, idOwner int64) (domain.Cart, error)
	AddCartItems(ctx context.Context, items []domain.Item, cardId int64, userId int64) error

	ShowCartItem(ctx context.Context, idOwner int64) (map[int64][]domain.Item, error)
}
