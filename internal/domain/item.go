package domain

import (
	"context"
)

type ItemRepo interface {
	ChangeName(ctx context.Context, cartID int64, oldItemName, newItemName string) error
	Delete(ctx context.Context, cartID int64, name string) error
}

type Item struct {
	Name   string `json:"name"`
	Amount uint8  `json:"amount"`
}
