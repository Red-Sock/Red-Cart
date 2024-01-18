package carts

import (
	"context"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/jackc/pgx/v5"

	"github.com/Red-Sock/Red-Cart/internal/clients/postgres"
	"github.com/Red-Sock/Red-Cart/internal/domain/cart"
	"github.com/Red-Sock/Red-Cart/internal/domain/user"
)

type Carts struct {
	conn postgres.Conn
}

func New(conn postgres.Conn) *Carts {
	return &Carts{conn: conn}
}

func (c *Carts) GetByCartId(ctx context.Context, cartId int64) (cart.Cart, error) {
	var dbCart cart.Cart

	err := c.conn.QueryRow(ctx, `
	SELECT 
	    id, 
	    owner_id 
	FROM cart
	WHERE id = $1`,
		cartId).
		Scan(
			&dbCart.Id,
			&dbCart.OwnerId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dbCart, nil
		}

		return cart.Cart{}, errors.Wrap(err, "error getting cart by ownerId from database")
	}
	return dbCart, nil
}

func (c *Carts) GetByOwnerId(ctx context.Context, ownerId int64) (cart.Cart, error) {
	var dbCart cart.Cart
	err := c.conn.QueryRow(ctx, `
	SELECT 
		id,
		owner_id
	FROM cart
	WHERE owner_id = $1`,
		ownerId).
		Scan(
			&dbCart.Id,
			&dbCart.OwnerId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dbCart, nil
		}

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

func (c *Carts) AddCartItems(ctx context.Context, items []string, cartId int64, userId int64) error {
	_, err := c.conn.Exec(ctx, `
	INSERT INTO carts_items (cart_id, item_name, user_id)
	VALUES 					(	  $1, 		 $2, 	  $3)
	ON CONFLICT (user_id, cart_id)
	DO UPDATE SET item_name = array_cat(carts_items.item_name, $4)`,
		cartId, items, userId, items)
	if err != nil {
		return errors.Wrap(err, "error add cartItem")
	}
	return nil
}

func (c *Carts) ShowCartItems(ctx context.Context, ownerId int64) ([]cart.CartItem, error) {
	var dbCart cart.Cart
	dbCart, err := c.GetByOwnerId(ctx, ownerId)
	if err != nil {
		return nil, errors.Wrap(err, "error show cart")
	}

	row, err := c.conn.Query(ctx, `
	SELECT 
		item_name,user_id
	FROM carts_items
	WHERE cart_id = $1`,
		dbCart.Id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "error getting cart by ownerId from database")
	}

	cartItem := make([]cart.CartItem, 0)
	for i := 0; row.Next(); i++ {
		item := cart.CartItem{}
		err = row.Scan(&item.ItemNames, &item.UserID)
		if err != nil {
			return nil, errors.Wrap(err, "error getting cart by ownerId from database")
		}
		cartItem = append(cartItem, item)
	}
	return cartItem, nil
}

func (c *Carts) GetUser(ctx context.Context, userId int64) (user.User, error) {
	var dbUser user.User
	err := c.conn.QueryRow(ctx, `
	SELECT 
		tg_id,
		user_name,
		first_name,
		last_name
	FROM tg_users
	WHERE tg_id = $1`,
		userId,
	).Scan(
		&dbUser.Id,
		&dbUser.UserName,
		&dbUser.FirstName,
		&dbUser.LastName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dbUser, nil
		}
		return user.User{}, errors.Wrap(err, "error getting user from database")
	}

	return dbUser, nil
}