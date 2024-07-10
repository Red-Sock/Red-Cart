package handlers

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
	"github.com/Red-Sock/Red-Cart/scripts"
)

func (d *DefaultHandler) basicInputs(msgIn *model.MessageIn, userCart domain.UserCart, out tgapi.Chat) (bool, error) {
	instruction, ok := d.expectedInstructions[scripts.GetLang(msgIn.From.LanguageCode)][msgIn.Text]
	if !ok {
		return false, nil
	}

	_ = out.SendMessage(&response.DeleteMessage{
		ChatId:    msgIn.Chat.ID,
		MessageId: int64(msgIn.MessageID),
	})
	var msg tgapi.MessageOut
	var err error
	switch instruction {
	case scripts.OpenSetting:
		msg, err = message.CartSettings(msgIn.Ctx, out, userCart)
	case scripts.Clear:
		msg, err = message.Delete(msgIn.Ctx, out, userCart)
	default:
		err = out.SendMessage(response.NewMessage(string("cannot handle " + instruction)))
		if err != nil {
			return false, errors.Wrap(err)
		}

		return true, nil
	}

	if err != nil {
		return true, errors.Wrap(err, "error assembling message")
	}

	if msg == nil {
		return true, nil
	}

	if !msgIn.IsCallback {
		_ = out.SendMessage(&response.DeleteMessage{
			ChatId:    msgIn.Chat.ID,
			MessageId: int64(msgIn.MessageID),
		})
	}

	err = d.cartService.SyncCartMessage(msgIn.Ctx, userCart, msg)
	if err != nil {
		err = out.SendMessage(response.NewMessage(err.Error()))
		if err != nil {
			return false, errors.Wrap(err)
		}

		return true, nil
	}

	return true, nil
}
