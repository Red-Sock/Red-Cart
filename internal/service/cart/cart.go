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

	if cart.OwnerId != 0 {
		return "", errors.New(fmt.Sprintf("У вас уже есть корзина с идентификатором = %d", cart.Id))
	}

	id, err := c.cartsData.Create(ctx, idOwner)
	if err != nil {
		return "", errors.Wrap(err, "error Creating cart")
	}

	return fmt.Sprintf(msgString, id, id), nil
}

func (c *CartsService) AddCartItems(ctx context.Context, items []domain.Item, cardId int64, userId int64) error {
	cartFromDB, err := c.cartsData.GetById(ctx, cardId)
	if err != nil {
		return err
	}

	if cartFromDB.Id == 0 {
		outMsg := fmt.Sprintf("Корзины с id = %d не существует", cardId)
		return errors.New(outMsg)
	}

	err = c.cartsData.AddCartItems(ctx, items, cardId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (c *CartsService) GetByCartId(ctx context.Context, cartId int64) (domain.Cart, error) {
	return c.cartsData.GetById(ctx, cartId)
}

func (c *CartsService) GetByOwnerId(ctx context.Context, idOwner int64) (*domain.Cart, error) {
	return c.cartsData.GetByOwnerId(ctx, idOwner)
}

func (c *CartsService) ShowCartItem(ctx context.Context, idOwner int64) (map[int64][]domain.Item, error) {
	return c.cartsData.ListCartItems(ctx, idOwner)
}
