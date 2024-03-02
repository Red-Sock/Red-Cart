package handlers

import (
	"encoding/json"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

func (d *DefaultHandler) editItemName(in *model.MessageIn, out tgapi.Chat, cart domain.UserCart) {
	if in.Text == "" {
		out.SendMessage(response.NewMessage("in order to change name you have to pass a valid string name"))
		return
	}
	var p domain.ChangeItemNamePayload
	err := json.Unmarshal(cart.Cart.StatePayload, &p)
	if err != nil {
		out.SendMessage(response.NewMessage("error parsing cart payload"))
		return
	}

	out.SendMessage(&response.DeleteMessage{
		ChatId:    in.Chat.ID,
		MessageId: int64(in.MessageID),
	})

	err = d.itemService.UpdateName(in.Ctx, cart.Cart.ID, p.ItemName, in.Text)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	err = d.cartService.AwaitItemsAdded(in.Ctx, cart.Cart.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	cart, err = d.cartService.GetCartById(in.Ctx, cart.Cart.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	msg := message.CartFromDomain(in.Ctx, out, cart)

	err = d.cartService.SyncCartMessage(in.Ctx, cart.Cart, msg)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}
}
