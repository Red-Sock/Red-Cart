package settings

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
	cartService service.CartService
}

func New(cartService service.CartService) *Handler {
	return &Handler{
		cartService: cartService,
	}
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) error {
	cart, err := h.cartService.GetCartByChatId(in.Ctx, in.Chat.ID)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	msg, err := message.CartSettings(in.Ctx, out, cart)
	if err != nil {
		return errors.Wrap(err, "error assembling cart settings message")
	}

	err = h.cartService.SyncCartMessage(in.Ctx, cart, msg)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	return nil
}

func (h *Handler) GetCommand() string {
	return commands.CartSetting
}
