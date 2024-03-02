package handlers

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/delete_cart"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/settings"
	"github.com/Red-Sock/Red-Cart/scripts"
)

func (d *DefaultHandler) basicInputs(in *model.MessageIn, out tgapi.Chat) bool {
	instruction, ok := d.expectedInstructions[scripts.GetLang(in.From.LanguageCode)][in.Text]
	if !ok {
		return false
	}

	var handler tgapi.Handler

	switch instruction {
	case scripts.OpenSetting:
		handler = settings.New(d.cartService)
	case scripts.Clear:
		handler = delete_cart.New(d.cartService)
		in.Text += " "
	default:
		out.SendMessage(response.NewMessage(string("cannot handle " + instruction)))
		return true
	}

	handler.Handle(in, out)

	return true
}
