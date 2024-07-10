package open_clear_menu

import (
	"strconv"

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

// Handle expects to have cart id as an argument
func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) error {
	if len(in.Args) < 1 {
		return out.SendMessage(response.NewMessage("required to have cart id as an argument"))
	}

	cartId, err := strconv.ParseInt(in.Args[0], 10, 64)
	if err != nil {
		return out.SendMessage(response.NewMessage("required to have cart id - integer"))
	}

	if !in.IsCallback {
		_ = out.SendMessage(&response.DeleteMessage{
			ChatId:    in.Chat.ID,
			MessageId: int64(in.MessageID),
		})
	}

	cart, err := h.cartService.GetCartByChatId(in.Ctx, cartId)
	if err != nil {
		return out.SendMessage(response.NewMessage("error getting cart for current chat"))
	}

	if len(cart.Cart.Items) == 0 {
		return nil
	}

	_, err = message.Delete(in.Ctx, out, cart)
	if err != nil {
		return errors.Wrap(err, "error assembling delete message")
	}

	return nil
}

func (h *Handler) GetCommand() string {
	return commands.ClearMenu
}
