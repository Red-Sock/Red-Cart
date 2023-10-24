package carts

import (
	"context"
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
}

func New() *Carts {
	return &Carts{
		idCartMap: make(map[int64]*cart.Cart),
		ownerMap:  make(map[int64]*cart.Cart),
	}
}

func (c *Carts) GetByCartId(ctx context.Context, cartId int64) (cart.Cart, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	cartNew, ok := c.idCartMap[cartId]
	if !ok {
		return cart.Cart{}, errors.New("Корзина не найдена")
	}
	return *cartNew, nil
}

func (c *Carts) GetByOwnerId(ctx context.Context, ownerId int64) (cart.Cart, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	cartNew, ok := c.ownerMap[ownerId]
	if !ok {
		return cart.Cart{}, errors.New("Корзина не найдена")

	}
	return *cartNew, nil
}

func (c *Carts) Create(ctx context.Context, idOwner int64) (id int64, err error) {
	idCartNew := c.idCart.Add(1)
	c.rw.Lock()
	defer c.rw.Unlock()

	newCart := &cart.Cart{
		Id:      idCartNew,
		OwnerId: idOwner,
		Url:     "",
	}
	c.idCartMap[idCartNew] = newCart
	c.ownerMap[idOwner] = newCart

	return idCartNew, nil
}

func (c *Carts) Show(ctx context.Context, id int64) (cart.Cart, error) {
	//TODO Чет я уже забыл где я это буду использовать, но обязательно буду!)
	return cart.Cart{}, errors.New("not implemented")
}
