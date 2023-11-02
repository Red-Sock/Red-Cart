package carts

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain/cart"
)

type Carts struct {
	idCart    atomic.Int64
	rw        sync.RWMutex
	idCartMap map[int64]*cart.Cart
	ownerMap  map[int64]*cart.Cart
	cartItems map[int64][]cart.CartItem
}

func New() *Carts {
	return &Carts{
		idCartMap: make(map[int64]*cart.Cart),
		ownerMap:  make(map[int64]*cart.Cart),
		cartItems: make(map[int64][]cart.CartItem),
	}
}

func (c *Carts) GetByCartId(ctx context.Context, cartId int64) (cart.Cart, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	cartNew, ok := c.idCartMap[cartId]
	if !ok {
		return cart.Cart{}, nil
	}

	return *cartNew, nil
}

func (c *Carts) GetByOwnerId(ctx context.Context, ownerId int64) (cart.Cart, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	result, ok := c.ownerMap[ownerId]
	if !ok {
		return cart.Cart{}, nil
	}

	return *result, nil
}

func (c *Carts) Create(ctx context.Context, idOwner int64) (id int64, err error) {
	newCart := &cart.Cart{
		Id:      c.idCart.Add(1),
		OwnerId: idOwner,
	}

	c.rw.Lock()
	defer c.rw.Unlock()

	c.idCartMap[newCart.Id] = newCart
	c.ownerMap[idOwner] = newCart

	return newCart.Id, nil
}

func (c *Carts) AddCartItems(ctx context.Context, items []string, cardId int64, userId int64) error {
	c.cartItems[cardId] = append(c.cartItems[cardId], cart.CartItem{
		ItemNames: items,
		UserID:    userId,
	})
	//TODO Пока что просто для проверки, что все отработало
	c.ShowCart(ctx, cardId)
	return nil
}

// TODO В дальнейшем будет выводить список товаров в корзине
func (c *Carts) ShowCart(ctx context.Context, cardId int64) (cart.Cart, error) {
	for _, ListItems := range c.cartItems[cardId] {
		fmt.Println(ListItems)
	}
	return cart.Cart{}, errors.New("not implemented")
}
