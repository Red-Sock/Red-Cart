package message

import (
	"context"
	"strconv"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/scripts"
)

func ClearCart(ctx context.Context, cart domain.UserCart) tgapi.MessageOut {
	if len(cart.Cart.Items) == 0 {
		return nil
	}

	cartIdStr := strconv.FormatUint(uint64(cart.Cart.ID), 10)

	keys := &keyboard.GridKeyboard{}
	keys.Columns = 1

	items, key := itemList(cart.Cart.Items)
	for i, itemName := range items {
		keys.AddButton(
			itemName+scripts.DeleteIcon,
			commands.DeleteItem+" "+cartIdStr+" "+key[i])
	}

	keys.AddButton(scripts.Get(ctx, scripts.PurgeCart), commands.Purge+" "+cartIdStr)
	keys.AddButton(scripts.BackIcon, commands.Cart)

	text := scripts.Get(ctx, scripts.ClickToRemove)

	if cart.Cart.MessageId != nil {
		return &response.EditMessage{
			ChatId:    cart.Cart.ChatId,
			Text:      text,
			MessageId: *cart.Cart.MessageId,
			Keys:      keys,
		}
	}

	msg := response.NewMessage(text)
	msg.AddKeyboard(keys)

	return msg
}
