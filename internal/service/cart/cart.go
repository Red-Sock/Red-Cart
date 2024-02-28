package cart

import (
	"context"
	"fmt"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

const msgString = `Корзина c id = %d была успешно создана.
Друзья могут добавить корзину через
/add_item %d имя_товара_1 имя_товара_2`

type CartsService struct {
	cartsData domain.CartRepo
}

func New(cartData domain.CartRepo) *CartsService {
	return &CartsService{
		cartsData: cartData,
	}
}

func (c *CartsService) Create(ctx context.Context, idOwner int64) (string, error) {
	cart, err := c.cartsData.GetByOwnerId(ctx, idOwner)
	if err != nil {
		return "", errors.New("Ошибка БД при получения корзины по Id")
	}

	if cart != nil {
		return "", errors.New(fmt.Sprintf("У вас уже есть корзина с идентификатором = %d", cart.ID))
	}

	cartId, err := c.cartsData.Create(ctx, idOwner)
	if err != nil {
		return "", errors.Wrap(err, "error creating cart")
	}

	err = c.cartsData.SetDefaultCart(ctx, idOwner, cartId)
	if err != nil {
		return "", errors.Wrap(err, "error setting default cart")
	}

	return fmt.Sprintf(msgString, cartId, cartId), nil
}

func (c *CartsService) UpdateMessageRef(ctx context.Context, cart domain.Cart) error {
	return c.cartsData.UpdateCart(ctx, cart)
}
