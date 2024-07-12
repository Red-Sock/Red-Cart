package message

import (
	"context"
	"strconv"

	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"

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
	keys := CartKeys(userCart.Cart)

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

func CartKeys(cart domain.Cart) (keys *keyboard.FloatingKeyboard) {
	if len(cart.Items) == 0 {
		return nil
	}
	keys = &keyboard.FloatingKeyboard{}

	items, itemKeys := itemList(cart.Items)
	for i := range items {
		row := make([]keyboard.Button, 0, 4)

		row = append(row, getCheckItemButton(cart.Items[i], itemKeys[i]))

		keys.AddRow(row)
	}

	return keys
}

func getCheckItemButton(item domain.Item, itemKey string) (key keyboard.Button) {

	key.Text = item.Name
	if item.Amount > 1 {
		key.Text += "(" + strconv.Itoa(int(item.Amount)) + ")"
	}

	if item.Checked {
		key.Text += scripts.CheckedIcon
		key.Value = commands.NewUncheckCommand(itemKey)
	} else {
		key.Value = commands.NewCheckCommand(itemKey)
	}

	return key
}

func getDeleteItemFromCart(itemKey string) (key keyboard.Button) {
	return keyboard.Button{
		Text:  scripts.BinIcon,
		Value: commands.NewDeleteCommand(itemKey),
	}
}
