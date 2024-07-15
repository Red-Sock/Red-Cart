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
	"github.com/Red-Sock/Red-Cart/scripts/icons"
)

func OpenCart(ctx context.Context, userCart domain.UserCart) interfaces.MessageOut {
	var text string

	if len(userCart.Cart.Items) == 0 {
		return emptyCart(ctx, userCart)
	}

	text = scripts.Get(ctx, scripts.Cart)

	keys := &keyboard.GridKeyboard{
		Columns: 1,
	}
	for _, item := range userCart.Cart.Items {
		keys.AddButton(getCheckItemButton(item))
	}

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

func getCheckItemButton(item domain.Item) (key keyboard.Button) {
	key.Text = getItemButtonName(item)
	if item.Checked {
		key.Value = commands.NewUncheckCommand(item.Name)
	} else {
		key.Value = commands.NewCheckCommand(item.Name)
	}

	return key
}

func getDeleteItemButton(item domain.Item) (key keyboard.Button) {
	key.Text = getItemButtonName(item)
	key.Value = commands.NewDeleteCommand(item.Name)

	return key
}
func getItemButtonName(item domain.Item) string {
	name := item.Name
	if item.Amount > 1 {
		name += "(" + strconv.Itoa(int(item.Amount)) + ")"
	}
	if item.Checked {
		name += icons.CheckedIcon
	}

	return name
}
