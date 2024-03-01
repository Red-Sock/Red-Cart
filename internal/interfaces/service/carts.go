package service

import (
	"context"

	tgapi "github.com/Red-Sock/go_tg/interfaces"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type CartService interface {
	SyncCartMessage(ctx context.Context, cart domain.Cart, msg tgapi.MessageOut) error

	GetCartByChatId(ctx context.Context, chatID int64) (domain.UserCart, error)

	Add(ctx context.Context, items []domain.Item, cartID int64, userID int64) (domain.UserCart, error)

	GetCartById(ctx context.Context, cartID int64) (domain.UserCart, error)

	AwaitNameChange(ctx context.Context, cartID int64, item domain.Item) error
	AwaitItemsAdded(ctx context.Context, cartID int64) error
}
