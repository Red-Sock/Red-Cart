package service

import (
	"context"
)

type ItemService interface {
	UpdateName(ctx context.Context, cartId int64, oldName, newName string) error
	Delete(ctx context.Context, cartId int64, itemName string) error
	Check(ctx context.Context, cartId int64, itemName string) error
	Uncheck(ctx context.Context, cartId int64, itemName string) error
}
