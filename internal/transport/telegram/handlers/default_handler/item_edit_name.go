package default_handler

import (
	"encoding/json"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

func (d *DefaultHandler) editItemName(msgIn *model.MessageIn, cart domain.UserCart) (tgapi.MessageOut, error) {
	if msgIn.Text == "" {
		return nil, errors.New("in order to change name you have to pass a valid string name")
	}

	var p domain.ChangeItemNamePayload
	err := json.Unmarshal(cart.Cart.StatePayload, &p)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	err = d.itemService.UpdateName(msgIn.Ctx, cart.Cart.ID, p.ItemName, msgIn.Text)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	err = d.cartService.AwaitItemsAdded(msgIn.Ctx, cart.Cart.ID)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	cart, err = d.cartService.GetCartById(msgIn.Ctx, cart.Cart.ID)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return message.OpenCart(msgIn.Ctx, cart), nil
}
