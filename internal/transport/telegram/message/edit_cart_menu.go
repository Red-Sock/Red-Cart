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

func EditFromCartItem(ctx context.Context, chat interfaces.Chat, userCart domain.UserCart, item domain.Item) (interfaces.MessageOut, error) {
	msgTxt := item.Name + "( " + strconv.Itoa(int(item.Amount)) + " ) "

	cartId := strconv.Itoa(int(userCart.Cart.ID))

	keys := keyboard.Keyboard{Columns: 2}
	keys.AddButton("🔙", commands.Cart)
	keys.AddButton(scripts.Get(ctx, scripts.Rename), commands.Rename+" "+cartId+" "+item.Name)

	if userCart.Cart.MessageId != nil {
		out := &response.EditMessage{
			ChatId:    userCart.Cart.ChatId,
			MessageId: *userCart.Cart.MessageId,
			Text:      msgTxt,
			Keys:      &keys,
		}

		err := chat.SendMessage(out)

		if err == nil {
			return out, nil
		}
	}

	msg := response.NewMessage(msgTxt)
	msg.AddKeyboard(keys)

	return msg, chat.SendMessage(msg)
}
