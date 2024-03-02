package message

import (
	"context"
	"sort"
	"strconv"

	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/scripts"
)

func CartFromDomain(ctx context.Context, chat interfaces.Chat, cart domain.UserCart) interfaces.MessageOut {
	var text string

	if len(cart.Cart.Items) == 0 {
		text = scripts.Get(ctx, scripts.CartIsEmpty)
		msg := response.NewMessage(text)
		chat.SendMessage(msg)

		return msg
	}

	text = "Корзина"

	var keys *keyboard.Keyboard

	if len(cart.Cart.Items) != 0 {
		sort.Slice(cart.Cart.Items, func(i, j int) bool {
			return cart.Cart.Items[i].Name < cart.Cart.Items[j].Name
		})
		keys = &keyboard.Keyboard{}
		keys.Columns = 1
		for _, item := range cart.Cart.Items {
			keys.AddButton(item.Name+" ( "+strconv.FormatUint(uint64(item.Amount), 10)+" )", commands.Check+" "+item.Name)
		}
	}

	if cart.Cart.MessageID != nil {
		out := &response.EditMessage{
			ChatId:    cart.Cart.ChatID,
			MessageId: *cart.Cart.MessageID,
			Text:      text,
			Keys:      keys,
		}
		chat.SendMessage(out)

		if out.GetMessageId() != 0 {
			return out
		}
	}

	out := &response.MessageOut{
		ChatId: cart.User.ID,
		Text:   text,
		Keys:   keys,
	}
	chat.SendMessage(out)
	return out
}

func CartSettings(chat interfaces.Chat, cart domain.UserCart) interfaces.MessageOut {
	var text string
	if len(cart.Cart.Items) == 0 {
		text = "Корзина пуста"
	} else {
		text = "Корзина"
	}

	var keys *keyboard.Keyboard

	if len(cart.Cart.Items) != 0 {
		keys = &keyboard.Keyboard{}
		keys.Columns = 1
		for _, item := range cart.Cart.Items {
			keys.AddButton(item.Name+" ( "+strconv.FormatUint(uint64(item.Amount), 10)+" )", commands.Edit+" "+item.Name)
		}

		keys.AddButton("️🔙", commands.Cart)
	}

	if cart.Cart.MessageID != nil {
		out := &response.EditMessage{
			ChatId:    cart.Cart.ChatID,
			MessageId: *cart.Cart.MessageID,
			Text:      text,
			Keys:      keys,
		}
		chat.SendMessage(out)

		if out.GetMessageId() != 0 {
			return out
		}
	}

	out := &response.MessageOut{
		ChatId: cart.User.ID,
		Text:   text,
		Keys:   keys,
	}
	chat.SendMessage(out)
	return out
}
