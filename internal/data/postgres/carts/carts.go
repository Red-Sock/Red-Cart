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
	    (owner_id, state)
	VALUES	(   $1,   $2)
	RETURNING id`,
		idOwner, domain.CartStateAdding).
		Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "error creating cart")
	}

	return id, nil
}

func (c *Repo) SetDefaultCart(ctx context.Context, userID int64, cartID int64) error {
	_, err := c.conn.Exec(ctx, `
	UPDATE carts_users
	SET is_default = (cart_id = $1)
	WHERE user_id = $2
`, cartID, userID)
	if err != nil {
		return errors.Wrap(err, "error setting default cart")
	}

	return nil
}

func (c *Repo) LinkUserToCart(ctx context.Context, userId int64, cartId int64) error {
	_, err := c.conn.Exec(ctx, `
	INSERT INTO carts_users
			(user_id, cart_id, is_default) 
	VALUES  ($1, $2, $3)
	ON CONFLICT DO NOTHING 
	`, userId, cartId, false)
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
			cu.chat_id,
			cu.message_id
		FROM carts_users cu
		LEFT JOIN public.carts c ON c.id = cu.cart_id
		WHERE cu.user_id = $1
		AND   cu.is_default`,
		userID).
		Scan(
			&cart.ID,
			&cart.ChatId,
			&cart.MessageId,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Cart{}, nil
		}

		return domain.Cart{}, errors.Wrap(err, "error getting user's default cart")
	}

	return cart, nil
}

func (c *Repo) GetByOwnerId(ctx context.Context, ownerId int64) (*domain.UserCart, error) {
	var dbCart domain.UserCart
	err := c.conn.QueryRow(ctx, `
	SELECT 
		c.id,
		c.owner_id,
		
		cu.chat_id,
		cu.message_id,
		
		c.state,
		c.state_payload
	FROM carts c
	JOIN carts_users cu ON cu.cart_id = c.id
	AND cu.user_id = $1
	WHERE owner_id = $1`,
		ownerId).
		Scan(
			&dbCart.Cart.ID,
			&dbCart.User.ID,
			&dbCart.Cart.ChatId,
			&dbCart.Cart.MessageId,
			&dbCart.Cart.State,
			&dbCart.Cart.StatePayload,
		)
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
		user_id,
		checked
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
		err = row.Scan(&item.Name, &item.Amount, &userID, &item.Checked)
		if err != nil {
			return nil, errors.Wrap(err, "error getting cart by ownerId from database")
		}

		cartItem = append(cartItem, item)
	}

	return cartItem, nil
}

func (c *Repo) UpdateCartReference(ctx context.Context, cart domain.UserCart) error {
	_, err := c.conn.Exec(ctx, `
		UPDATE carts_users
		SET
		    chat_id = $1,
		    message_id = $2 
		WHERE cart_id = $3 AND user_id = $4
`, cart.Cart.ChatId, cart.Cart.MessageId, cart.Cart.ID, cart.User.ID)
	if err != nil {
		return errors.Wrap(err, "error updating cart")
	}

	return nil
}

func (c *Repo) GetCartByChatId(ctx context.Context, chatId int64) (resp *domain.UserCart, err error) {
	resp = &domain.UserCart{}

	err = c.conn.QueryRow(ctx, `
		SELECT 
		    c.id,
		    c.owner_id,
		    
			cu.chat_id,
			cu.message_id,
			c.state,
			c.state_payload
		FROM carts c
		LEFT JOIN carts_users cu 
		ON cu.cart_id = c.id 
		AND cu.chat_id = $1
`, chatId).
		Scan(
			&resp.Cart.ID,
			&resp.User.ID,

			&resp.Cart.ChatId,
			&resp.Cart.MessageId,
			&resp.Cart.State,
			&resp.Cart.StatePayload,
		)
	if err != nil {
		return resp, errors.Wrap(err, "error scanning cart")
	}

	return resp, nil
}

func (c *Repo) GetCartByID(ctx context.Context, id int64) (*domain.UserCart, error) {
	resp := &domain.UserCart{}
	err := c.conn.QueryRow(ctx, `
		SELECT 
		    c.id,
		    c.owner_id,

			cu.chat_id,
			cu.message_id,

			c.state,
			c.state_payload
		FROM carts c
		LEFT JOIN carts_users cu ON c.id = cu.cart_id
		WHERE id = $1
`, id).
		Scan(
			&resp.Cart.ID,
			&resp.User.ID,

			&resp.Cart.ChatId,
			&resp.Cart.MessageId,

			&resp.Cart.State,
			&resp.Cart.StatePayload,
		)
	if err != nil {
		return nil, errors.Wrap(err, "error getting cart")
	}

	return resp, nil
}

func (c *Repo) ChangeState(ctx context.Context, req domain.Cart) error {
	_, err := c.conn.Exec(ctx, `
		UPDATE carts 
		SET 
		    state = $1,
		    state_payload = $2 
		WHERE id = $3
`, req.State, req.StatePayload, req.ID)

	if err != nil {
		return errors.Wrap(err, "error updating cart state")
	}

	return nil
}

func (c *Repo) PurgeCart(ctx context.Context, cartId int64) error {
	_, err := c.conn.Exec(ctx, `
		DELETE 
		FROM cart_items
	   	WHERE cart_id = $1
`, cartId)
	if err != nil {
		return errors.Wrap(err, "error purging items in cart")
	}

	return nil
}
