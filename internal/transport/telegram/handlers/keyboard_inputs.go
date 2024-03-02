package handlers

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
	"github.com/Red-Sock/Red-Cart/scripts"
)

func (d *DefaultHandler) basicInputs(in *model.MessageIn, userCart domain.UserCart, out tgapi.Chat) (bool, error) {
	instruction, ok := d.expectedInstructions[scripts.GetLang(in.From.LanguageCode)][in.Text]
	if !ok {
		return false, nil
	}

	_ = out.SendMessage(&response.DeleteMessage{
		ChatId:    in.Chat.ID,
		MessageId: int64(in.MessageID),
	})
	var msg tgapi.MessageOut
	var err error
	switch instruction {
	case scripts.OpenSetting:
		msg, err = message.CartSettings(in.Ctx, out, userCart)
	case scripts.Clear:
		msg, err = message.Delete(in.Ctx, out, userCart)
	default:
		return true, out.SendMessage(response.NewMessage(string("cannot handle " + instruction)))
	}

	if err != nil {
		return true, err
	}

	if msg == nil {
		return true, nil
	}

	if !in.IsCallback {
		_ = out.SendMessage(&response.DeleteMessage{
			ChatId:    in.Chat.ID,
			MessageId: int64(in.MessageID),
		})
	}

	err = d.cartService.SyncCartMessage(in.Ctx, userCart, msg)
	if err != nil {
		return true, out.SendMessage(response.NewMessage(err.Error()))
	}

	return true, nil
}
