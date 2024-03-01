package message

import (
	"strconv"

	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
)

func EditFromCartItem(chat interfaces.Chat, userCart domain.UserCart, item domain.Item) interfaces.MessageOut {
	msgTxt := item.Name + "( " + strconv.Itoa(int(item.Amount)) + " ) "

	keys := keyboard.InlineKeyboard{Columns: 2}
	keys.AddButton("🔙", commands.Cart)
	keys.AddButton("Rename✏️", commands.Rename+" "+strconv.Itoa(int(userCart.Cart.ID))+" "+item.Name)

	if userCart.Cart.MessageID != nil {
		out := &response.EditMessage{
			ChatId:    userCart.Cart.ChatID,
			MessageId: *userCart.Cart.MessageID,
			Text:      msgTxt,
			Keys:      &keys,
		}

		chat.SendMessage(out)

		if out.GetMessageId() != 0 {
			return out
		}
	}

	msg := response.NewMessage(msgTxt)
	msg.AddKeyboard(keys)

	chat.SendMessage(msg)

	return msg
}
