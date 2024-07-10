package message

import (
	"context"
	"strconv"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/scripts"
)

// nolint
func Delete(ctx context.Context, out tgapi.Chat, cart domain.UserCart) (tgapi.MessageOut, error) {
	if len(cart.Cart.Items) == 0 {
		return nil, nil
	}

	cartIdStr := strconv.FormatUint(uint64(cart.Cart.ID), 10)

	keys := keyboard.Keyboard{}
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
		msg := &response.EditMessage{
			ChatId:    cart.Cart.ChatId,
			Text:      text,
			MessageId: *cart.Cart.MessageId,
			Keys:      &keys,
		}

		err := out.SendMessage(msg)
		if err == nil {
			return msg, errors.Wrap(err, "error sending message")
		}
	}

	msg := response.NewMessage(text)
	msg.AddKeyboard(keys)
	err := out.SendMessage(msg)
	if err != nil {
		return nil, errors.Wrap(err, "error sending message")
	}

	return msg, nil
}
