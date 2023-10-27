package carts

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain/cart"
)

type CartItem struct {
	m map[int64]*cart.CartItem
}

func NewCartItem() *CartItem {
	return &CartItem{
		m: make(map[int64]*cart.CartItem),
	}
}

func (c *CartItem) CreateCartsItem(ctx context.Context, cartId int64) (id int64, err error) {
	return 0, nil
}
