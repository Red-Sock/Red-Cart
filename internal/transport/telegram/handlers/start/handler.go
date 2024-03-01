package start

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
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

	out.SendMessage(response.NewMessage(startMessage.Msg))

	out.SendMessage(&response.DeleteMessage{
		ChatId:    in.Chat.ID,
		MessageId: int64(in.MessageID),
	})

	if startMessage.Cart.MessageID != nil {
		out.SendMessage(&response.DeleteMessage{
			ChatId:    startMessage.Cart.ChatID,
			MessageId: *startMessage.Cart.MessageID,
		})

		startMessage.Cart.MessageID = nil
	}

	cartMsg := message.CartFromDomain(out, startMessage.UserCart)

	err = h.cartSrv.SyncCartMessage(in.Ctx, startMessage.Cart, cartMsg)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}
}

func (h *Handler) GetDescription() string {
	return "returns just hello"
}

func (h *Handler) GetCommand() string {
	return commands.Start
}
