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
	cartItem  *CartItemService
}

type CartItemService struct {
	cartsData data.Carts
}

func NewCartItem() *CartItemService {
	return &CartItemService{}
}

func New(userData data.Carts) *CartsService {
	return &CartsService{
		cartsData: userData,
		cartItem:  NewCartItem(),
	}
}

func (c *CartsService) Create(ctx context.Context, idOwner int64) (string, error) {
	cart, err := c.cartsData.GetByOwnerId(ctx, idOwner)

	//TODO пока что не знаю что делать с этой ошибкой, так как в нормальной БД мы можем оставить эту ошибку, а для карты она нам не нужна

	//if err != nil {
	//	msg := fmt.Sprintf("Ошибка БД при получения корзины по Id")
	//	return "", errors.New(msg)
	//}

	if cart.OwnerId != 0 {
		msg := fmt.Sprintf("У вас уже есть корзина с идентификатором = %d", cart.Id)
		return "", errors.New(msg)
	}

	id, err := c.cartsData.Create(ctx, idOwner)

	if err != nil {
		return "", errors.Wrap(err, "error Creating cart")
	}

	return fmt.Sprintf(msgString, id, id), nil
}

func (c *CartsService) AddCartItems(ctx context.Context, items []string, cardId int64, userId int64) error {
	_, err := c.GetByCartId(ctx, cardId)
	if err != nil {
		outMsg := fmt.Sprintf("Корзины с id = %d  не существует", cardId)
		return errors.New(outMsg)
	}
	//TODO добавить логику с ошибкой и возвратом ответа, если он нужен
	c.cartsData.AddCartItems(ctx, items, cardId, userId)

	return nil
}

func (c *CartsService) GetByCartId(ctx context.Context, cartId int64) (cart.Cart, error) {
	cart, err := c.cartsData.GetByCartId(ctx, cartId)
	if err != nil {
		return cart, err
	}

	return cart, nil
}

func (c *CartsService) GetByOwnerId(ctx context.Context, idOwner int64) (cart.Cart, error) {
	cart, err := c.cartsData.GetByOwnerId(ctx, idOwner)
	if err != nil {
		return cart, err
	}

	return cart, nil
}
