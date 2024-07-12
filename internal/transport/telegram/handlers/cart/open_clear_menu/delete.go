package open_clear_menu

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
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

// Handle expects to have cart id as an argument
func (h *Handler) Handle(msgIn *model.MessageIn, out tgapi.Chat) error {
	defer helpers.DeleteIncomingMessage(msgIn, out)

	cart, err := h.cartService.GetCartByChatId(msgIn.Ctx, msgIn.Chat.ID)
	if err != nil {
		return out.SendMessage(response.NewMessage("error getting cart for current chat"))
	}

	if len(cart.Cart.Items) == 0 {
		return nil
	}

	msg := message.ClearCartMenu(msgIn.Ctx, cart)
	err = out.SendMessage(msg)
	if err != nil {
		return errors.Wrap(err, "error assembling delete message")
	}

	return nil
}

func (h *Handler) GetCommand() string {
	return commands.ClearMenu
}
