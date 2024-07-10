package domain

import (
	"context"
)

type ItemRepo interface {
	ChangeName(ctx context.Context, cartId int64, oldItemName, newItemName string) error
	Delete(ctx context.Context, cartId int64, itemName string) error
	Check(ctx context.Context, cartId int64, itemName string) error
	Uncheck(ctx context.Context, cartId int64, itemName string) error
}

type Item struct {
	Name    string `json:"name"`
	Amount  uint8  `json:"amount"`
	Checked bool   `json:"checked"`
}
