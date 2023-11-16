package cart

import (
	"context"
	"fmt"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain/cart"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
)

const msgString = `Корзина c id = %d была успешно создана.
Друзья могут добавить корзину через
/add_item %d имя_товара_1 имя_товара_2`

type CartsService struct {
	cartsData data.Carts
}

func New(cartData data.Carts) *CartsService {
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

func (c *CartsService) AddCartItems(ctx context.Context, items []string, cardId int64, userId int64) error {
	cartFromDB, err := c.cartsData.GetByCartId(ctx, cardId)
	if err != nil {
		return err
	}

	if cartFromDB.Id == 0 {
		outMsg := fmt.Sprintf("Корзины с id = %d не существует", cardId)
		return errors.New(outMsg)
	}
	//TODO [RC-12] добавить логику с ошибкой и возвратом ответа, если он нужен
	c.cartsData.AddCartItems(ctx, items, cardId, userId)

	return nil
}

func (c *CartsService) GetByCartId(ctx context.Context, cartId int64) (cart.Cart, error) {
	res, err := c.cartsData.GetByCartId(ctx, cartId)
	if err != nil {
		return cart.Cart{}, err
	}

	return res, nil
}

func (c *CartsService) GetByOwnerId(ctx context.Context, idOwner int64) (cart.Cart, error) {
	res, err := c.cartsData.GetByOwnerId(ctx, idOwner)
	if err != nil {
		return cart.Cart{}, err
	}

	return res, nil
}
