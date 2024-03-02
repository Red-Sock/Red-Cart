package handlers

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

func (d *DefaultHandler) addItem(in *model.MessageIn, out tgapi.Chat, userCart domain.UserCart) {
	items := make([]domain.Item, 0, len(in.Args))
	for _, item := range in.Args {
		items = append(items, domain.Item{Name: item, Amount: 1})
	}

	cart, err := d.cartService.Add(in.Ctx, items, userCart.Cart.ID, in.From.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	out.SendMessage(&response.DeleteMessage{
		ChatId:    in.Chat.ID,
		MessageId: int64(in.MessageID),
	})

	msg := message.CartFromDomain(in.Ctx, out, cart)

	err = d.cartService.SyncCartMessage(in.Ctx, userCart.Cart, msg)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}
}
