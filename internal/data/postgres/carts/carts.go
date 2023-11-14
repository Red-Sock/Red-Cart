package carts

import (
	"context"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/jackc/pgx/v5"

	"github.com/Red-Sock/Red-Cart/internal/clients/postgres"
	"github.com/Red-Sock/Red-Cart/internal/domain/cart"
)

type Carts struct {
	conn postgres.Conn
}

func New(conn postgres.Conn) *Carts {
	return &Carts{conn: conn}
}

func (c *Carts) GetByCartId(ctx context.Context, cartId int64) (cart.Cart, error) {
	var dbCart cart.Cart

	err := c.conn.QueryRow(ctx, `SELECT id, owner_id FROM cart
	    WHERE id = $1`, cartId).Scan(&dbCart.Id, &dbCart.OwnerId)
	if errors.Is(err, pgx.ErrNoRows) {
		return dbCart, nil
	}

	if err != nil {
		return cart.Cart{}, errors.Wrap(err, "error getting cart by ownerId from database")
	}
	return dbCart, nil
}

func (c *Carts) GetByOwnerId(ctx context.Context, ownerId int64) (cart.Cart, error) {
	var dbCart cart.Cart
	err := c.conn.QueryRow(ctx, `SELECT id, owner_id FROM cart
	    WHERE owner_id = $1`, ownerId).Scan(&dbCart.Id, &dbCart.OwnerId)
	if errors.Is(err, pgx.ErrNoRows) {
		return dbCart, nil
	}

	if err != nil {
		return cart.Cart{}, errors.Wrap(err, "error getting cart by ownerId from database")
	}
	return dbCart, nil
}

func (c *Carts) Create(ctx context.Context, idOwner int64) (id int64, err error) {
	_, err = c.conn.Exec(ctx, `
INSERT INTO cart
	    (owner_id)
VALUES	(   $1)`,
		idOwner,
	)
	if err != nil {
		return 0, errors.Wrap(err, "error creating cart")
	}

	dbCart, err := c.GetByOwnerId(ctx, idOwner)
	if err != nil {
		return 0, errors.Wrap(err, "error creating cart")
	}

	return dbCart.Id, nil
}

func (c *Carts) AddCartItems(ctx context.Context, items []string, cardId int64, userId int64) error {
	_, err := c.conn.Exec(ctx, `INSERT INTO carts_items
	(item_name,user_id)
	VALUES($1,$2)`,
		items, userId)

	if err != nil {
		return errors.Wrap(err, "error add cartItem")
	}
	return nil
}

// TODO В дальнейшем будет выводить список товаров в корзине
func (c *Carts) ShowCart(ctx context.Context, cardId int64) (cart.Cart, error) {
	return cart.Cart{}, errors.New("not implemented")
}
