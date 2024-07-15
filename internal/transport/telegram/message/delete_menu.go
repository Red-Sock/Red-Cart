package message

import (
	"context"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/scripts"
	"github.com/Red-Sock/Red-Cart/scripts/icons"
)

func ClearCartMenu(ctx context.Context, userCart domain.UserCart) tgapi.MessageOut {
	if len(userCart.Cart.Items) == 0 {
		return emptyCart(ctx, userCart)
	}

	keys := &keyboard.GridKeyboard{}
	keys.Columns = 1

	for _, item := range userCart.Cart.Items {
		keys.AddButton(getDeleteItemButton(item))
	}

	keys.AddButton(keyboard.NewButton(scripts.Get(ctx, scripts.PurgeCartAction), commands.Purge))
	keys.AddButton(keyboard.NewButton(icons.BackIcon, commands.Cart))
	text := icons.BinIcon

	if userCart.Cart.MessageId != nil {
		return &response.EditMessage{
			ChatId:    userCart.Cart.ChatId,
			Text:      text,
			MessageId: *userCart.Cart.MessageId,
			Keys:      keys,
		}
	}

	msg := response.NewMessage(text)
	msg.AddKeyboard(keys)

	return msg
}
