package start

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	errors "github.com/Red-Sock/trace-errors"
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
	defer func() {
		err := out.SendMessage(&response.DeleteMessage{
			ChatId:    msgIn.Chat.ID,
			MessageId: int64(msgIn.MessageID),
		})
		if err != nil {
			logrus.Errorf("error sending delete message command: %s", err)
		}
	}()

	newUser := parsing.ToDomainUser(msgIn)

	startMessage, err := h.userSrv.Start(msgIn.Ctx, newUser, msgIn.Chat.ID)
	if err != nil {
		return errors.Wrap(err)
	}

	h.removePreviousMessage(&startMessage, out)

	msg := response.NewMessage(startMessage.Msg)
	err = out.SendMessage(msg)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
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
