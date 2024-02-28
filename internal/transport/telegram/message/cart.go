package message

import (
	"strconv"

	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
)

func CartFromDomain(cart domain.UserCart) interfaces.MessageOut {
	var text string
	if len(cart.Cart.Items) == 0 {
		text = "Корзина пуста"
	} else {
		text = "Корзина"
	}

	var keys *keyboard.InlineKeyboard

	if len(cart.Cart.Items) != 0 {
		keys = &keyboard.InlineKeyboard{}
		keys.Columns = 3
		for _, item := range cart.Cart.Items {
			keys.AddButton(item.Name+"-"+strconv.FormatUint(uint64(item.Amount), 10), "/start")
			keys.AddButton("+", "/start")
			keys.AddButton("-", "/start")
		}
	}

	if cart.Cart.MessageID != nil && cart.Cart.ChatID != nil {
		return &response.EditMessage{
			ChatId:    *cart.Cart.ChatID,
			MessageId: *cart.Cart.MessageID,
			Text:      text,
			Keys:      keys,
		}
	}

	return &response.MessageOut{
		ChatId: cart.User.ID,
		Text:   text,
		Keys:   keys,
	}
}
