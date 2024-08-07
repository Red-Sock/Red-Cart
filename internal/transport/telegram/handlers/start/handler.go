package start

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"
	errors "github.com/Red-Sock/trace-errors"
	"github.com/sirupsen/logrus"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/helpers"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/parsing"
	"github.com/Red-Sock/Red-Cart/scripts"
)

type Handler struct {
	userSrv service.UserService
	cartSrv service.CartService
}

func New(srv service.Service) *Handler {
	return &Handler{
		userSrv: srv.User(),
		cartSrv: srv.Cart(),
	}
}

func (h *Handler) Handle(msgIn *model.MessageIn, out tgapi.Chat) error {
	defer helpers.DeleteIncomingMessage(msgIn, out)

	newUser := parsing.ToDomainUser(msgIn)

	startMessage, err := h.userSrv.Start(msgIn.Ctx, newUser, msgIn.Chat.ID)
	if err != nil {
		return errors.Wrap(err)
	}

	if startMessage.Cart.MessageId != nil {
		h.removePreviousMessage(&startMessage, out)
	}

	startMsg := response.NewMessage(startMessage.Msg)

	startKeyboard := &keyboard.GridKeyboard{}
	startKeyboard.AddButton(keyboard.NewButton(scripts.Get(msgIn.Ctx, scripts.OpenClearMenu), ""))
	startKeyboard.SetIsReplyKeyboard(true)
	startMsg.Keys = startKeyboard

	err = out.SendMessage(startMsg)
	if err != nil {
		return errors.Wrap(err)
	}

	cartMsg := message.OpenCart(msgIn.Ctx, startMessage.UserCart)
	err = out.SendMessage(cartMsg)
	if err != nil {
		return errors.Wrap(err)
	}

	err = h.cartSrv.SyncCartMessage(msgIn.Ctx, startMessage.UserCart, cartMsg)
	if err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (h *Handler) removePreviousMessage(startMessage *domain.StartMessagePayload, out tgapi.Chat) {
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
