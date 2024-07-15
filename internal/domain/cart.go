package domain

import (
	"context"
)

const (
	CartStateAdding          = "adding"
	CartStateEditingItemName = "editing_item_name"
)

type CartRepo interface {
	Create(ctx context.Context, idOwner int64) (id int64, err error)

	SetDefaultCart(ctx context.Context, userID, cartID int64) error
	LinkUserToCart(ctx context.Context, userID, cartID, chatId int64) error

	ListCartItems(ctx context.Context, cartId int64) ([]Item, error)

	AddCartItems(ctx context.Context, items []Item, cardId int64, userId int64) error

	UpdateCartReference(ctx context.Context, cart UserCart) error

	GetUserDefaultCart(ctx context.Context, id int64) (Cart, error)
	GetByOwnerId(ctx context.Context, id int64) (*UserCart, error)

	GetCartByChatId(ctx context.Context, chatID int64) (*UserCart, error)
	GetCartByID(ctx context.Context, id int64) (*UserCart, error)

	ChangeState(ctx context.Context, req Cart) error
	PurgeCart(ctx context.Context, cartId int64) error
}

type CartFilter struct {
	CartId  []int64
	OwnerID []int64

	Paging
}

type cartState string

type Cart struct {
	Id    int64
	Items []Item

	ChatId    int64
	MessageId *int64

	State        cartState
	StatePayload []byte
}

type ChangeItemNamePayload struct {
	ItemName string `json:"item_name"`
}
