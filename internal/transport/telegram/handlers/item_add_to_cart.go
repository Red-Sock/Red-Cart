package handlers

import (
	"strings"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

func (d *DefaultHandler) addItem(in *model.MessageIn, out tgapi.Chat, userCart domain.UserCart) error {
	itemsRaw := strings.Split(in.Text, "\n")
	items := make([]domain.Item, 0, len(itemsRaw))

	for _, item := range itemsRaw {
		items = append(items, domain.Item{Name: item, Amount: 1})
	}

	cart, err := d.cartService.Add(in.Ctx, items, userCart.Cart.ID, in.From.ID)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	_ = out.SendMessage(&response.DeleteMessage{
		ChatId:    in.Chat.ID,
		MessageId: int64(in.MessageID),
	})

	msg, err := message.OpenCart(in.Ctx, out, cart)
	if err != nil {
		return err
	}

	err = d.cartService.SyncCartMessage(in.Ctx, userCart.Cart, msg)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	return nil
}
