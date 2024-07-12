package default_handler

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
	"github.com/Red-Sock/Red-Cart/scripts"
)

func (d *DefaultHandler) basicInputs(msgIn *model.MessageIn, userCart domain.UserCart) tgapi.MessageOut {
	lang := scripts.GetLang(msgIn.From.LanguageCode)
	langInstr, ok := d.expectedInstructions[lang]
	if !ok {
		return nil
	}

	instruction, ok := langInstr[msgIn.Text]
	if !ok {
		return nil
	}

	switch instruction {
	case scripts.OpenSetting:
		return message.CartSettings(msgIn.Ctx, userCart)
	case scripts.Clear:
		return message.ClearCart(msgIn.Ctx, userCart)
	default:
		return response.NewMessage(string("cannot handle " + instruction))
	}
}
