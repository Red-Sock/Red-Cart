package message

import (
	"context"
	"strconv"

	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/scripts"
)

// nolint
func OpenCart(ctx context.Context, chat interfaces.Chat, cart domain.UserCart) (interfaces.MessageOut, error) {
	var text string

	//nolint
	if len(cart.Cart.Items) == 0 {
		text = scripts.Get(ctx, scripts.CartIsEmpty)

		var msg interfaces.MessageOut
		if cart.Cart.MessageId != nil {
			msg = &response.EditMessage{
				ChatId:    cart.Cart.ChatId,
				Text:      text,
				MessageId: *cart.Cart.MessageId,
			}
		} else {
			msg = response.NewMessage(text)
		}
		err := chat.SendMessage(msg)
		if err != nil {
			return nil, errors.Wrap(err, "error sending cart message")
		}

		return msg, nil
	}

	text = scripts.Get(ctx, scripts.Cart)

	var keys *keyboard.Keyboard
	cartId := strconv.Itoa(int(cart.Cart.ID))
	if len(cart.Cart.Items) != 0 {
		keys = &keyboard.Keyboard{}
		keys.Columns = 1

		items, itemKeys := itemList(cart.Cart.Items)
		for i, itemName := range items {
			if !cart.Cart.Items[i].Checked {
				keys.AddButton(itemName, commands.Check+" "+cartId+" "+itemKeys[i])
			} else {
				keys.AddButton(itemName+" "+scripts.CheckedIcon, commands.Uncheck+" "+cartId+" "+itemKeys[i])
			}
		}
	}

	if cart.Cart.MessageId != nil {
		out := &response.EditMessage{
			ChatId:    cart.Cart.ChatId,
			MessageId: *cart.Cart.MessageId,
			Text:      text,
			Keys:      keys,
		}
		err := chat.SendMessage(out)
		if err == nil {
			return out, nil
		}
	}

	out := &response.MessageOut{
		ChatId: cart.User.ID,
		Text:   text,
		Keys:   keys,
	}
	err := chat.SendMessage(out)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return out, nil
}

// nolint
func CartSettings(ctx context.Context, chat interfaces.Chat, cart domain.UserCart) (interfaces.MessageOut, error) {
	var text string
	if len(cart.Cart.Items) == 0 {
		text = scripts.Get(ctx, scripts.CartIsEmpty)
	} else {
		text = scripts.Get(ctx, scripts.Cart)
	}

	var keys *keyboard.Keyboard

	if len(cart.Cart.Items) != 0 {
		keys = &keyboard.Keyboard{}
		keys.Columns = 1

		itemsNames, itemKeys := itemList(cart.Cart.Items)
		for i, itemName := range itemsNames {
			keys.AddButton(itemName, commands.Edit+" "+itemKeys[i])
		}

		keys.AddButton("️🔙", commands.Cart)
	}

	if cart.Cart.MessageId != nil {
		out := &response.EditMessage{
			ChatId:    cart.Cart.ChatId,
			MessageId: *cart.Cart.MessageId,
			Text:      text,
			Keys:      keys,
		}
		err := chat.SendMessage(out)

		if err == nil {
			return out, nil
		}
	}

	out := &response.MessageOut{
		ChatId: cart.User.ID,
		Text:   text,
		Keys:   keys,
	}

	err := chat.SendMessage(out)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return out, nil
}
