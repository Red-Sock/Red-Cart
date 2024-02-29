package cart

import (
	"context"
	"fmt"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

const msgString = `Корзина c id = %d была успешно создана.
Друзья могут добавить корзину через
/add_item %d имя_товара_1 имя_товара_2`

type CartsService struct {
	cartData domain.CartRepo
}

func New(cartData domain.CartRepo) *CartsService {
	return &CartsService{
		cartData: cartData,
	}
}

func (c *CartsService) Create(ctx context.Context, idOwner int64) (string, error) {
	cart, err := c.cartData.GetByOwnerId(ctx, idOwner)
	if err != nil {
		return "", errors.New("Ошибка БД при получения корзины по Id")
	}

	if cart != nil {
		return "", errors.New(fmt.Sprintf("У вас уже есть корзина с идентификатором = %d", cart.ID))
	}

	cartId, err := c.cartData.Create(ctx, idOwner)
	if err != nil {
		return "", errors.Wrap(err, "error creating cart")
	}

	err = c.cartData.SetDefaultCart(ctx, idOwner, cartId)
	if err != nil {
		return "", errors.Wrap(err, "error setting default cart")
	}

	return fmt.Sprintf(msgString, cartId, cartId), nil
}

func (c *CartsService) SyncCartMessage(ctx context.Context, cart domain.Cart, msg tgapi.MessageOut) error {
	if cart.MessageID == nil && cart.ChatID == nil {
		return nil
	}

	chatID, msgID := msg.GetChatId(), msg.GetMessageId()
	cart.MessageID = &msgID
	cart.ChatID = &chatID

	err := c.cartData.UpdateCartReference(ctx, cart)
	if err != nil {
		return errors.Wrap(err, "error updating cart chat reference")
	}

	return nil
}
