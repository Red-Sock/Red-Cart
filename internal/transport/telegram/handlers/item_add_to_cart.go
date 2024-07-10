package handlers

import (
	"strings"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

func (d *DefaultHandler) addItem(msgIn *model.MessageIn, out tgapi.Chat, userCart domain.UserCart) error {
	itemsRaw := strings.Split(msgIn.Text, "\n")
	items := make([]domain.Item, 0, len(itemsRaw))

	for _, item := range itemsRaw {
		items = append(items, domain.Item{Name: item, Amount: 1})
	}

	cart, err := d.cartService.Add(msgIn.Ctx, items, userCart.Cart.ID, msgIn.From.ID)
	if err != nil {
		err = out.SendMessage(response.NewMessage(err.Error()))
		if err != nil {
			return errors.Wrap(err)
		}

		return nil
	}

	_ = out.SendMessage(&response.DeleteMessage{
		ChatId:    msgIn.Chat.ID,
		MessageId: int64(msgIn.MessageID),
	})

	msg, err := message.OpenCart(msgIn.Ctx, out, cart)
	if err != nil {
		return errors.Wrap(err, "error assembling open cart message")
	}

	err = d.cartService.SyncCartMessage(msgIn.Ctx, userCart, msg)
	if err != nil {
		err = out.SendMessage(response.NewMessage(err.Error()))
		if err != nil {
			return errors.Wrap(err)
		}

		return nil
	}

	return nil
}
