package cart

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

const Command = "/cart"

type Handler struct {
	userService service.UserService
	cartService service.CartService
}

func New(userService service.UserService, cartService service.CartService) *Handler {
	return &Handler{
		userService: userService,
		cartService: cartService,
	}
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {
	userCart, err := h.userService.GetDefaultCart(in.Ctx, in.From.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	out.SendMessage(&response.DeleteMessage{
		ChatId:    in.Chat.ID,
		MessageId: int64(in.MessageID),
	})

	msg := message.CartFromDomain(userCart)
	out.SendMessage(msg)

	if msg.GetMessageId() == 0 {
		userCart.Cart.MessageID = nil
		msg = message.CartFromDomain(userCart)
		out.SendMessage(msg)
	}

	err = h.cartService.SyncCartMessage(in.Ctx, userCart.Cart, msg)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}
}

func (h *Handler) GetDescription() string {
	return "Shows default cart"
}

func (h *Handler) GetCommand() string {
	return Command
}
