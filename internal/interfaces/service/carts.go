package service

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain/cart"
)

type CartService interface {
	Create(ctx context.Context, idOwner int64) (string, error)
	GetByOwnerId(ctx context.Context, idOwner int64) (cart.Cart, error)
	GetByCartId(ctx context.Context, idOwner int64) (cart.Cart, error)
	AddCartItems(ctx context.Context, items []string, cardId int64, userId int64) error
	ShowCartItem(ctx context.Context, idOwner int64) ([]cart.CartItem, error)
}
