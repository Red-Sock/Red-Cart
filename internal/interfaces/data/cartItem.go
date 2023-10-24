package data

import (
	"context"
)

type CartsItem interface {
	//TODO пока просто функция заглушка для структуры проекта
	CreateCartsItem(ctx context.Context, cartId int64) (id int64, err error)
}
