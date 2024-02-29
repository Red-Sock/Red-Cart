package domain

import (
	"context"
)

type CartRepo interface {
	Create(ctx context.Context, idOwner int64) (id int64, err error)
	SetDefaultCart(ctx context.Context, userID int64, cartID int64) error
	LinkUserToCart(ctx context.Context, userID int64, cartID int64) error

	ListCartItems(ctx context.Context, cartId int64) ([]Item, error)

	AddCartItems(ctx context.Context, items []Item, cardId int64, userId int64) error

	UpdateCartReference(ctx context.Context, cart Cart) error

	GetUserDefaultCart(ctx context.Context, id int64) (Cart, error)
	GetByOwnerId(ctx context.Context, id int64) (*Cart, error)
}

type CartFilter struct {
	CartId  []int64
	OwnerID []int64

	Paging
}

type Cart struct {
	ID    int64
	Items []Item

	ChatID    *int64
	MessageID *int64
}

type Item struct {
	Name   string `json:"name"`
	Amount uint8  `json:"amount"`
}
