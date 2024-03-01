package message

import (
	"strconv"

	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
)

func CartFromDomain(chat interfaces.Chat, cart domain.UserCart) interfaces.MessageOut {
	var text string
	if len(cart.Cart.Items) == 0 {
		text = "Корзина пуста"
	} else {
		text = "Корзина"
	}

	var keys *keyboard.InlineKeyboard
	cartId := strconv.Itoa(int(cart.Cart.ChatID))

	if len(cart.Cart.Items) != 0 {
		keys = &keyboard.InlineKeyboard{}
		keys.Columns = 1
		for _, item := range cart.Cart.Items {
			keys.AddButton(item.Name+" ( "+strconv.FormatUint(uint64(item.Amount), 10)+" )", commands.Edit+" "+item.Name)
		}

		keys.AddButton("Удалить товар / очистить корзину 🗑️", commands.Delete+" "+cartId)
	}

	if cart.Cart.MessageID != nil {
		out := &response.EditMessage{
			ChatId:    cart.Cart.ChatID,
			MessageId: *cart.Cart.MessageID,
			Text:      text,
			Keys:      keys,
		}
		chat.SendMessage(out)

		if out.GetMessageId() != 0 {
			return out
		}
	}

	out := &response.MessageOut{
		ChatId: cart.User.ID,
		Text:   text,
		Keys:   keys,
	}
	chat.SendMessage(out)
	return out
}
