package carts

import (
	"context"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/jackc/pgx/v5"

	"github.com/Red-Sock/Red-Cart/internal/clients/postgres"
	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type Carts struct {
	conn postgres.Conn
}

func New(conn postgres.Conn) *Carts {
	return &Carts{conn: conn}
}

func (c *Carts) GetById(ctx context.Context, cartId int64) (domain.Cart, error) {
	var dbCart domain.Cart

	err := c.conn.QueryRow(ctx, `
	SELECT 
	    id, 
	    owner_id 
	FROM carts
	WHERE id = $1`,
		cartId).
		Scan(
			&dbCart.Id,
			&dbCart.OwnerId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dbCart, nil
		}

		return domain.Cart{}, errors.Wrap(err, "error getting cart by ownerId from database")
	}
	return dbCart, nil
}

func (c *Carts) GetByOwnerId(ctx context.Context, ownerId int64) (domain.Cart, error) {
	var dbCart domain.Cart
	err := c.conn.QueryRow(ctx, `
	SELECT 
		id,
		owner_id
	FROM carts
	WHERE owner_id = $1`,
		ownerId).
		Scan(
			&dbCart.Id,
			&dbCart.OwnerId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dbCart, nil
		}

		return domain.Cart{}, errors.Wrap(err, "error getting cart by ownerId from database")
	}
	return dbCart, nil
}

func (c *Carts) Create(ctx context.Context, idOwner int64) (id int64, err error) {
	_, err = c.conn.Exec(ctx, `
	INSERT INTO carts
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

func (c *Carts) AddCartItems(ctx context.Context, items []domain.Item, cartId int64, userId int64) error {
	_, err := c.conn.CopyFrom(ctx,
		[]string{"cart_items"},
		[]string{"cart_id", "item_name", "amount", "user_id"},
		pgx.CopyFromSlice(len(items), func(i int) ([]any, error) {
			if i >= len(items) {
				return nil, pgx.ErrTooManyRows
			}
			return []any{cartId, items[i].Name, items[i].Amount, userId}, nil
		}),
	)
	if err != nil {
		return errors.Wrap(err, "error add cartItem")
	}

	return nil
}

func (c *Carts) ListCartItems(ctx context.Context, ownerId int64) (map[int64][]domain.Item, error) {
	var dbCart domain.Cart
	dbCart, err := c.GetByOwnerId(ctx, ownerId)
	if err != nil {
		return nil, errors.Wrap(err, "error show cart")
	}

	row, err := c.conn.Query(ctx, `
	SELECT 
		item_name,user_id
	FROM cart_items
	WHERE cart_id = $1`,
		dbCart.Id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "error getting cart by ownerId from database")
	}

	cartItem := make(map[int64][]domain.Item)
	for row.Next() {
		item := domain.Item{}
		var userID int64
		err = row.Scan(&item.Name, &userID)
		if err != nil {
			return nil, errors.Wrap(err, "error getting cart by ownerId from database")
		}

		cartItem[userID] = append(cartItem[userID], item)
	}
	return cartItem, nil
}

func (c *Carts) GetUser(ctx context.Context, userId int64) (domain.User, error) {
	var dbUser domain.User
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
		return domain.User{}, errors.Wrap(err, "error getting user from database")
	}

	return dbUser, nil
}
