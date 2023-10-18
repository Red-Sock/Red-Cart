package carts

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain/cart"
)

var idCart atomic.Int64

type Carts struct {
	rw sync.RWMutex
	m  map[int64]cart.Cart
}

func New() *Carts {
	return &Carts{
		m: make(map[int64]cart.Cart),
	}
}

func (c Carts) Get(ctx context.Context, idOwner int64) (cart.Cart, error) {

	cart, ok := c.m[idOwner]

	//TODO переделать логику Get
	if ok {
		//TODO не знаю правильно ли выкидывать ошибку пользователю
		msg := fmt.Sprintf("У вас уже есть корзина с идентификатором = %d", cart.Id)
		return c.m[idOwner], errors.New(msg)
	}
	return c.m[idOwner], nil
}

func (c Carts) Create(ctx context.Context, idOwner int64) (id int64, err error) {
	idCatNew := idCart.Add(1)
	c.rw.Lock()
	defer c.rw.Unlock()
	c.m[idOwner] = cart.Cart{
		Id:      idCatNew,
		OwnerId: idOwner,
		Url:     "",
	}
	return idCatNew, nil
}

func (c Carts) Show(ctx context.Context, id int64) (cart.Cart, error) {
	//TODO в рамках RC-11
	return cart.Cart{}, errors.New("not implemented")
}
