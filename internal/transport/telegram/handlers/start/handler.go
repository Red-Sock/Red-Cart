package start

import (
	"strconv"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"

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

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {
	newUser := domain.User{
		ID:        in.From.ID,
		UserName:  in.From.UserName,
		FirstName: in.From.FirstName,
		LastName:  in.From.LastName,
	}

	startMessage, err := h.userSrv.Start(in.Ctx, newUser, in.Chat.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	if startMessage.Cart.MessageID != nil {
		out.SendMessage(&response.DeleteMessage{
			ChatId:    startMessage.Cart.ChatID,
			MessageId: *startMessage.Cart.MessageID,
		})

		startMessage.Cart.MessageID = nil
	}

	h.startMessage(in, startMessage, out)

	cartMsg := message.CartFromDomain(in.Ctx, out, startMessage.UserCart)

	err = h.cartSrv.SyncCartMessage(in.Ctx, startMessage.Cart, cartMsg)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	out.SendMessage(&response.DeleteMessage{
		ChatId:    in.Chat.ID,
		MessageId: int64(in.MessageID),
	})

}

func (h *Handler) startMessage(in *model.MessageIn, payload domain.StartMessagePayload, out tgapi.Chat) {
	cartId := strconv.Itoa(int(payload.Cart.ID))

	msg := response.NewMessage(payload.Msg)

	keyboardReply := keyboard.Keyboard{IsReplyKeyboard: true}
	keyboardReply.AddButton(scripts.Get(in.Ctx, scripts.OpenSetting), commands.CartSetting+" "+cartId)
	keyboardReply.AddButton(scripts.Get(in.Ctx, scripts.Clear), commands.Delete+" "+cartId)

	msg.AddKeyboard(keyboardReply)

	out.SendMessage(msg)
}

func (h *Handler) GetDescription() string {
	return "returns just hello"
}

func (h *Handler) GetCommand() string {
	return commands.Start
}
