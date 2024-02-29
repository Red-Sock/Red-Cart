package service

import (
	"context"

	tgapi "github.com/Red-Sock/go_tg/interfaces"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type CartService interface {
	Create(ctx context.Context, idOwner int64) (string, error)

	SyncCartMessage(ctx context.Context, cart domain.Cart, msg tgapi.MessageOut) error
}
