package service

import (
	"context"
)

type CartItem interface {
	//TODO пока просто функция заглушка для структуры проекта
	Create(ctx context.Context, idCart int64) (string, error)
}
