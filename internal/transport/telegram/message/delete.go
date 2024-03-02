package message

import (
	"strconv"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
)

func Delete(out tgapi.Chat, cart domain.UserCart) tgapi.MessageOut {
	cartIdStr := strconv.FormatUint(uint64(cart.Cart.ID), 10)

	keys := keyboard.Keyboard{}
	keys.Columns = 1
	for _, item := range cart.Cart.Items {
		amountStr := strconv.FormatUint(uint64(item.Amount), 10)

		keys.AddButton(
			item.Name+" ( "+amountStr+" )"+"❌",
			commands.DeleteItem+" "+cartIdStr+" "+item.Name)
	}

	keys.AddButton("Очистить корзину🚮", commands.Purge+" "+cartIdStr)
	keys.AddButton("🔙", commands.Cart)

	text := "Нажмите для удаления"

	if cart.Cart.MessageID != nil {
		msg := &response.EditMessage{
			ChatId:    cart.Cart.ChatID,
			Text:      text,
			MessageId: *cart.Cart.MessageID,
			Keys:      &keys,
		}

		out.SendMessage(msg)

		if msg.GetMessageId() != 0 {
			return msg
		}
	}

	msg := response.NewMessage(text)
	msg.AddKeyboard(keys)

	out.SendMessage(msg)

	return msg
}
