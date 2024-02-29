package carts

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	errors "github.com/Red-Sock/trace-errors"
	"github.com/jackc/pgx/v5"

	"github.com/Red-Sock/Red-Cart/internal/clients/postgres"
	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type Repo struct {
	conn postgres.Conn
}

func New(conn postgres.Conn) *Repo {
	return &Repo{conn: conn}
}

func (c *Repo) Create(ctx context.Context, idOwner int64) (id int64, err error) {
	err = c.conn.QueryRow(ctx, `
	INSERT INTO carts
	    (owner_id)
	VALUES	(   $1)
	
	RETURNING id`,
		idOwner,
	).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "error creating cart")
	}

	return id, nil
}

func (c *Repo) SetDefaultCart(ctx context.Context, userID int64, cartID int64) error {
	_, err := c.conn.Exec(ctx, `
	UPDATE carts_users SET is_default = (cart_id = $1) WHERE user_id = $2
`, cartID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Repo) LinkUserToCart(ctx context.Context, userID int64, cartID int64) error {
	_, err := c.conn.Exec(ctx, `
	INSERT INTO carts_users
			(user_id, cart_id, is_default) 
	VALUES  ($1, $2, $3)
	`, userID, cartID, false)
	if err != nil {
		return errors.Wrap(err, "error executing db query")
	}

	return err
}

func (c *Repo) GetUserDefaultCart(ctx context.Context, userID int64) (domain.Cart, error) {
	var cart domain.Cart
	err := c.conn.QueryRow(ctx, `
		SELECT
			cu.cart_id,
			c.chat_id,
			c.message_id
		FROM carts_users cu
		LEFT JOIN public.carts c ON c.id = cu.cart_id
		WHERE cu.user_id = $1
		AND   cu.is_default`,
		userID).
		Scan(
			&cart.ID,
			&cart.ChatID,
			&cart.MessageID,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Cart{}, nil
		}

		return domain.Cart{}, errors.Wrap(err, "error getting user's default cart")
	}

	return cart, nil
}

func (c *Repo) GetByOwnerId(ctx context.Context, ownerId int64) (*domain.Cart, error) {
	var dbCart domain.Cart
	err := c.conn.QueryRow(ctx, `
	SELECT 
		id
	FROM carts
	WHERE owner_id = $1`,
		ownerId).
		Scan(
			&dbCart.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "error getting cart by ownerId from database")
	}
	return &dbCart, nil
}

func (c *Repo) AddCartItems(ctx context.Context, items []domain.Item, cartId int64, userId int64) error {
	q := sq.Insert("cart_items").
		Columns(
			"cart_id",
			"item_name",
			"amount",
			"user_id",
		).
		Suffix(`
ON CONFLICT ( cart_id, item_name, user_id)
DO UPDATE SET 
amount = cart_items.amount+excluded.amount`).
		PlaceholderFormat(sq.Dollar)

	for _, item := range items {
		q = q.Values(cartId, item.Name, item.Amount, userId)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error assembling sql")
	}

	_, err = c.conn.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "error add cartItem")
	}

	return nil
}

func (c *Repo) ListCartItems(ctx context.Context, id int64) ([]domain.Item, error) {
	row, err := c.conn.Query(ctx, `
	SELECT 
		item_name,
		amount,
		user_id
	FROM cart_items
	WHERE cart_id = $1`,
		id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "error getting cart by ownerId from database")
	}

	var cartItem []domain.Item
	for row.Next() {
		item := domain.Item{}
		var userID int64
		err = row.Scan(&item.Name, &item.Amount, &userID)
		if err != nil {
			return nil, errors.Wrap(err, "error getting cart by ownerId from database")
		}

		cartItem = append(cartItem, item)
	}

	return cartItem, nil
}

func (c *Repo) ListCartItemsWithRequesters(ctx context.Context, ownerId int64) (map[int64][]domain.Item, error) {
	var dbCart *domain.Cart
	dbCart, err := c.GetByOwnerId(ctx, ownerId)
	if err != nil {
		return nil, errors.Wrap(err, "error show cart")
	}

	row, err := c.conn.Query(ctx, `
	SELECT 
		item_name, 
		user_id
	FROM cart_items
	WHERE cart_id = $1`,
		dbCart.ID)
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

func (c *Repo) GetUser(ctx context.Context, userId int64) (domain.User, error) {
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
		&dbUser.ID,
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

func (c *Repo) UpdateCartReference(ctx context.Context, cart domain.Cart) error {
	_, err := c.conn.Exec(ctx, `
		UPDATE carts
		SET
		    chat_id = $1,
		    message_id = $2 
		WHERE id = $3
`, cart.ChatID, cart.MessageID, cart.ID)
	if err != nil {
		return errors.Wrap(err, "error updating cart")
	}

	return nil
}
