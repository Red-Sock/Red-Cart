package start

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	"github.com/sirupsen/logrus"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/parsing"
)

type Handler struct {
	userSrv service.UserService
	cartSrv service.CartService
}

func New(userSrv service.UserService, cartSrv service.CartService) *Handler {
	return &Handler{
		userSrv: userSrv,
		cartSrv: cartSrv,
	}
}

func (h *Handler) Handle(msgIn *model.MessageIn, out tgapi.Chat) error {
	newUser := parsing.ToDomainUser(msgIn)

	startMessage, err := h.userSrv.Start(msgIn.Ctx, newUser, msgIn.Chat.ID)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	h.removePreviousMessage(&startMessage, out)

	h.sendStartMessage(msgIn, startMessage, out)

	return out.SendMessage(&response.DeleteMessage{
		ChatId:    msgIn.Chat.ID,
		MessageId: int64(msgIn.MessageID),
	})
}

func (h *Handler) sendStartMessage(in *model.MessageIn, payload domain.StartMessagePayload, out tgapi.Chat) {
	//cartId := strconv.Itoa(int(payload.Cart.ID))

	msg := response.NewMessage(payload.Msg)

	//keyboardReply := keyboard.Keyboard{}
	//keyboardReply.AddButton(scripts.Get(in.Ctx, scripts.Cart), commands.Cart)
	//keyboardReply.AddButton(scripts.Get(in.Ctx, scripts.OpenSetting), commands.CartSetting+" "+cartId)
	//keyboardReply.AddButton(scripts.Get(in.Ctx, scripts.Clear), commands.ClearMenu+" "+cartId)
	//
	//msg.AddKeyboard(keyboardReply)

	err := out.SendMessage(msg)
	if err != nil {
		logrus.Error(err.Error())
	}
}

func (h *Handler) removePreviousMessage(startMessage *domain.StartMessagePayload, out tgapi.Chat) {
	if startMessage.Cart.MessageId == nil {
		return
	}

	err := out.SendMessage(&response.DeleteMessage{
		ChatId:    startMessage.Cart.ChatId,
		MessageId: *startMessage.Cart.MessageId,
	})
	if err != nil {
		logrus.Errorf("error deleting previous message (chatID = %d, messageId = %d, %s",
			startMessage.Cart.ChatId,
			*startMessage.Cart.MessageId,
			err.Error())
		return
	}

	startMessage.Cart.MessageId = nil

}

func (h *Handler) GetDescription() string {
	return "returns just hello"
}

func (h *Handler) GetCommand() string {
	return commands.Start
}
