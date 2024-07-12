package handlers

import (
	"encoding/json"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

func (d *DefaultHandler) editItemName(msgIn *model.MessageIn, out tgapi.Chat, cart domain.UserCart) error {
	if msgIn.Text == "" {
		err := out.SendMessage(response.NewMessage("in order to change name you have to pass a valid string name"))
		if err != nil {
			return errors.Wrap(err)
		}

		return nil
	}

	var p domain.ChangeItemNamePayload
	err := json.Unmarshal(cart.Cart.StatePayload, &p)
	if err != nil {
		err = out.SendMessage(response.NewMessage("error parsing cart payload"))
		if err != nil {
			return errors.Wrap(err)
		}

		return nil
	}

	_ = out.SendMessage(&response.DeleteMessage{
		ChatId:    msgIn.Chat.ID,
		MessageId: int64(msgIn.MessageID),
	})

	err = d.itemService.UpdateName(msgIn.Ctx, cart.Cart.ID, p.ItemName, msgIn.Text)
	if err != nil {
		err = out.SendMessage(response.NewMessage(err.Error()))
		if err != nil {
			return errors.Wrap(err)
		}

		return nil
	}

	err = d.cartService.AwaitItemsAdded(msgIn.Ctx, cart.Cart.ID)
	if err != nil {
		err = out.SendMessage(response.NewMessage(err.Error()))
		if err != nil {
			return errors.Wrap(err)
		}

		return nil
	}

	cart, err = d.cartService.GetCartById(msgIn.Ctx, cart.Cart.ID)
	if err != nil {
		err = out.SendMessage(response.NewMessage(err.Error()))
		if err != nil {
			return errors.Wrap(err)
		}

		return nil
	}

	msg := message.OpenCart(msgIn.Ctx, cart)
	err = out.SendMessage(msg)
	if err != nil {
		return errors.Wrap(err)
	}

	err = d.cartService.SyncCartMessage(msgIn.Ctx, cart, msg)
	if err != nil {
		err = out.SendMessage(response.NewMessage(err.Error()))
		if err != nil {
			return errors.Wrap(err)
		}

		return nil
	}

	return nil
}
