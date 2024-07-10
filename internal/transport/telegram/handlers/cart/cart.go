package cart

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

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

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) error {
	userCart, err := h.userService.GetCartByChat(in.Ctx, in.From.ID)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	if !in.IsCallback {
		_ = out.SendMessage(&response.DeleteMessage{
			ChatId:    in.Chat.ID,
			MessageId: int64(in.MessageID),
		})
	}

	msg, err := message.OpenCart(in.Ctx, out, userCart)
	if err != nil {
		return errors.Wrap(err, "error assembling open cart message")
	}

	err = h.cartService.SyncCartMessage(in.Ctx, userCart, msg)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}
	return nil
}

func (h *Handler) GetDescription() string {
	return "Shows default cart"
}

func (h *Handler) GetCommand() string {
	return commands.Cart
}
