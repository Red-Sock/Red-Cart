package default_handler

import (
	"strings"

	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

func (d *DefaultHandler) addItem(msgIn *model.MessageIn, userCart domain.UserCart) (interfaces.MessageOut, error) {
	itemsRaw := strings.Split(msgIn.Text, "\n")
	items := make([]domain.Item, 0, len(itemsRaw))

	for _, item := range itemsRaw {
		items = append(items, domain.Item{Name: item, Amount: 1})
	}

	cart, err := d.cartService.Add(msgIn.Ctx, items, userCart.Cart.Id, msgIn.From.ID)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return message.OpenCart(msgIn.Ctx, cart), nil
}
