package carts

import (
	"context"
	"fmt"
	"sync/atomic"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain/cart"
)

var idCart int64

type Carts struct {
	m map[int64]cart.Cart
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
	atomic.AddInt64(&idCart, 1)
	c.m[idOwner] = cart.Cart{
		Id:      idCart,
		OwnerId: idOwner,
		Url:     "",
	}
	return idCart, nil
}

func (c Carts) Show(ctx context.Context, id int64) (cart.Cart, error) {
	//TODO в рамках RC-11
	panic("implement me")
}
