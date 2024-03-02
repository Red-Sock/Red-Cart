package settings

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

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

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {
	cart, err := h.cartService.GetCartByChatId(in.Ctx, in.Chat.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	err = h.cartService.SyncCartMessage(in.Ctx, cart.Cart, message.CartSettings(out, cart))
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}
}

func (h *Handler) GetCommand() string {
	return commands.CartSetting
}
