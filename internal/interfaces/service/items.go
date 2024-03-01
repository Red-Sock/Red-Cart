package service

import (
	"context"
)

type ItemService interface {
	UpdateName(ctx context.Context, cartID int64, oldName, newName string) error
}
