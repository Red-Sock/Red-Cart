package domain

import (
	"context"
)

type CartRepo interface {
	Create(ctx context.Context, idOwner int64) (id int64, err error)

	ListCartItems(ctx context.Context, ownerId int64) (map[int64][]Item, error)

	GetByOwnerId(ctx context.Context, id int64) (*Cart, error)
	GetById(ctx context.Context, id int64) (Cart, error)
	AddCartItems(ctx context.Context, items []Item, cardId int64, userId int64) error
}

type Cart struct {
	Id           int64
	OwnerId      int64
	UsersToItems map[int64]Item
}

type Item struct {
	Name   string
	Amount uint8
}
