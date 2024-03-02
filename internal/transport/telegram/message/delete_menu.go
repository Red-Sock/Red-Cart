package message

import (
	"strconv"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
)

func Delete(out tgapi.Chat, cart domain.UserCart) (tgapi.MessageOut, error) {
	if len(cart.Cart.Items) == 0 {
		return nil, nil
	}

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

	if cart.Cart.MessageId != nil {
		msg := &response.EditMessage{
			ChatId:    cart.Cart.ChatId,
			Text:      text,
			MessageId: *cart.Cart.MessageId,
			Keys:      &keys,
		}

		err := out.SendMessage(msg)

		if err == nil {
			return msg, err
		}
	}

	msg := response.NewMessage(text)
	msg.AddKeyboard(keys)

	return msg, out.SendMessage(msg)
}
