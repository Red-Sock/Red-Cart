package service

import (
	"context"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

type UserService interface {
	Start(ctx context.Context, user domain.User, chatID int64) (message domain.StartMessagePayload, err error)

	AddToDefaultCart(ctx context.Context, items []domain.Item, userID int64) (domain.UserCart, error)
	GetCartByChat(ctx context.Context, chatID int64) (domain.UserCart, error)
}
