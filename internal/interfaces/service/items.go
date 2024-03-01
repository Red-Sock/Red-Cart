package service

import (
	"context"
)

type ItemService interface {
	UpdateName(ctx context.Context, cartID int64, oldName, newName string) error
	Delete(ctx context.Context, cartID int64, itemName string) error
}
