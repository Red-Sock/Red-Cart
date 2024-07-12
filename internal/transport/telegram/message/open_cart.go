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

func OpenCart(ctx context.Context, userCart domain.UserCart) interfaces.MessageOut {
	var text string

	if len(userCart.Cart.Items) == 0 {
		return emptyCart(ctx, userCart)
	}

	text = scripts.Get(ctx, scripts.Cart)
	keys := cartKeys(userCart.Cart)

	if userCart.Cart.MessageId != nil {
		return &response.EditMessage{
			ChatId:    userCart.Cart.ChatId,
			MessageId: *userCart.Cart.MessageId,
			Text:      text,
			Keys:      keys,
		}
	}

	return &response.MessageOut{
		ChatId: userCart.User.Id,
		Text:   text,
		Keys:   keys,
	}
}

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
		ChatId: cart.User.Id,
		Text:   text,
		Keys:   keys,
	}

	err := chat.SendMessage(out)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return out, nil
}

func emptyCart(ctx context.Context, cart domain.UserCart) interfaces.MessageOut {
	text := scripts.Get(ctx, scripts.CartIsEmpty)

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

	return msg
}

func cartKeys(cart domain.Cart) (keys *keyboard.Keyboard) {
	cartId := strconv.Itoa(int(cart.ID))
	if len(cart.Items) == 0 {
		return nil
	}
	keys = &keyboard.Keyboard{}
	keys.Columns = 1

	items, itemKeys := itemList(cart.Items)
	for i, itemName := range items {
		if !cart.Items[i].Checked {
			keys.AddButton(itemName, commands.Check+" "+cartId+" "+itemKeys[i])
		} else {
			keys.AddButton(itemName+" "+scripts.CheckedIcon, commands.Uncheck+" "+cartId+" "+itemKeys[i])
		}
	}

	return keys
}
