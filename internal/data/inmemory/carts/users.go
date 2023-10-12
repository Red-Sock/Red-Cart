package carts

import (
	"context"
	"github.com/Red-Sock/Red-Cart/internal/domain/cart"
	"github.com/Red-Sock/Red-Cart/internal/domain/user"
)

type Carts struct {
	m map[int64]user.User
}

func NewCarts() *Carts {
	return &Carts{
		m: make(map[int64]user.User),
	}
}

func (c2 Carts) Create(ctx context.Context, c cart.Cart) error {
	//TODO implement me
	panic("implement me")
}

func (c2 Carts) Show(ctx context.Context, id int64) (cart.Cart, error) {
	//TODO implement me
	panic("implement me")
}
