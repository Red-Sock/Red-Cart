package data

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain/cart"
)

type Carts interface {
	Create(ctx context.Context, idOwner int64) (id int64, err error)
	Show(ctx context.Context, id int64) (cart.Cart, error)
	GetByOwnerId(ctx context.Context, id int64) (cart.Cart, error)
	GetByCartId(ctx context.Context, id int64) (cart.Cart, error)
}
