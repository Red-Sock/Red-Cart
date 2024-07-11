package start

import (
	"strconv"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"
	errors "github.com/Red-Sock/trace-errors"
	"github.com/sirupsen/logrus"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
	"github.com/Red-Sock/Red-Cart/scripts"
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
	newUser := domain.User{
		Id:        msgIn.From.ID,
		UserName:  msgIn.From.UserName,
		FirstName: msgIn.From.FirstName,
		LastName:  msgIn.From.LastName,
	}

	startMessage, err := h.userSrv.Start(msgIn.Ctx, newUser, msgIn.Chat.ID)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	if startMessage.Cart.MessageId != nil {
		_ = out.SendMessage(&response.DeleteMessage{
			ChatId:    startMessage.Cart.ChatId,
			MessageId: *startMessage.Cart.MessageId,
		})

		startMessage.Cart.MessageId = nil
	}

	h.startMessage(msgIn, startMessage, out)

	cartMsg, err := message.OpenCart(msgIn.Ctx, out, startMessage.UserCart)
	if err != nil {
		return errors.Wrap(err, "open cart error")
	}

	err = h.cartSrv.SyncCartMessage(msgIn.Ctx, startMessage.UserCart, cartMsg)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	return out.SendMessage(&response.DeleteMessage{
		ChatId:    msgIn.Chat.ID,
		MessageId: int64(msgIn.MessageID),
	})
}

func (h *Handler) startMessage(in *model.MessageIn, payload domain.StartMessagePayload, out tgapi.Chat) {
	cartId := strconv.Itoa(int(payload.Cart.ID))

	msg := response.NewMessage(payload.Msg)

	keyboardReply := keyboard.Keyboard{IsReplyKeyboard: true}
	keyboardReply.AddButton(scripts.Get(in.Ctx, scripts.OpenSetting), commands.CartSetting+" "+cartId)
	keyboardReply.AddButton(scripts.Get(in.Ctx, scripts.Clear), commands.ClearMenu+" "+cartId)

	msg.AddKeyboard(keyboardReply)

	err := out.SendMessage(msg)
	if err != nil {
		logrus.Error(err.Error())
	}
}

func (h *Handler) GetDescription() string {
	return "returns just hello"
}

func (h *Handler) GetCommand() string {
	return commands.Start
}
