package handlers

import (
	"encoding/json"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

func (d *DefaultHandler) editItemName(in *model.MessageIn, out tgapi.Chat, cart domain.UserCart) error {
	if in.Text == "" {
		return out.SendMessage(response.NewMessage("in order to change name you have to pass a valid string name"))
	}

	var p domain.ChangeItemNamePayload
	err := json.Unmarshal(cart.Cart.StatePayload, &p)
	if err != nil {
		return out.SendMessage(response.NewMessage("error parsing cart payload"))
	}

	_ = out.SendMessage(&response.DeleteMessage{
		ChatId:    in.Chat.ID,
		MessageId: int64(in.MessageID),
	})

	err = d.itemService.UpdateName(in.Ctx, cart.Cart.ID, p.ItemName, in.Text)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	err = d.cartService.AwaitItemsAdded(in.Ctx, cart.Cart.ID)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	cart, err = d.cartService.GetCartById(in.Ctx, cart.Cart.ID)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	msg, err := message.OpenCart(in.Ctx, out, cart)
	if err != nil {
		return err
	}

	err = d.cartService.SyncCartMessage(in.Ctx, cart.Cart, msg)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	return nil
}
