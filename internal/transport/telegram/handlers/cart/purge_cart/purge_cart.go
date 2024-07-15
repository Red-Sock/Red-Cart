package purge_cart

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/helpers"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

type Handler struct {
	cartService service.CartService
}

func New(srv service.Service) *Handler {
	return &Handler{
		cartService: srv.Cart(),
	}
}

// Handle expects to have cart id as a given parameter
func (h *Handler) Handle(msgIn *model.MessageIn, out tgapi.Chat) error {
	defer helpers.DeleteIncomingMessage(msgIn, out)

	userCart, err := h.cartService.GetCartByChatId(msgIn.Ctx, msgIn.Chat.ID)
	if err != nil {
		return errors.Wrap(err)
	}

	err = h.cartService.PurgeCart(msgIn.Ctx, userCart.Cart.Id)
	if err != nil {
		return errors.Wrap(err)
	}

	userCart, err = h.cartService.GetCartById(msgIn.Ctx, userCart.Cart.Id)
	if err != nil {
		return errors.Wrap(err)
	}

	msg := message.OpenCart(msgIn.Ctx, userCart)
	err = out.SendMessage(msg)
	if err != nil {
		return errors.Wrap(err)
	}

	err = h.cartService.SyncCartMessage(msgIn.Ctx, userCart, msg)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

// GetCommand - возвращает комманду
func (h *Handler) GetCommand() string {
	return commands.Purge
}
