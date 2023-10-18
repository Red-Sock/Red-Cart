package cart

import (
	"context"
	"fmt"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
)

const msgString = `Корзина c id = %d была успешно создана.
Друзья могут добавить корзину через
/add_item %d имя_товара_1 имя_товара_2`

type CartsService struct {
	cartsData data.Carts
}

func New(userData data.Carts) *CartsService {
	return &CartsService{
		cartsData: userData,
	}
}

func (c CartsService) Create(ctx context.Context, idOwner int64) (string, error) {
	_, err := c.cartsData.Get(ctx, idOwner)
	if err != nil {
		return "", err
	}

	return nil, ""
}

func (c CartsService) Get(ctx context.Context, idOwner int64) (error, string) {
	//TODO переделать обработку ошибки
	_, err := c.cartsData.Get(ctx, idOwner)
	if err != nil {
		return err, ""
	}

	id, err := c.cartsData.Create(ctx, idOwner)

	if err != nil {
		return "", errors.Wrap(err, "error Creating cart")
	}

	return fmt.Sprintf(msgString, id, id), nil
}
