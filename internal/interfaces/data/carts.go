package data

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain/cart"
)

type Carts interface {
	Create(ctx context.Context, c cart.Cart) error
	Show(ctx context.Context, id int64) (cart.Cart, error)
}
