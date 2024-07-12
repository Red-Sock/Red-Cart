package purge_cart

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

func New(srv service.Service) *Handler {
	return &Handler{
		cartService: srv.Cart(),
	}
}

// Handle expects to have cart id as a given parameter
func (h *Handler) Handle(msgIn *model.MessageIn, out tgapi.Chat) error {
	if len(msgIn.Args) < 1 {
		return out.SendMessage(response.NewMessage("expecting to have a cart id as an argument"))
	}

	cartId, err := strconv.ParseInt(msgIn.Args[0], 10, 64)
	if err != nil {
		return errors.Wrap(err)
	}

	err = h.cartService.PurgeCart(msgIn.Ctx, cartId)
	if err != nil {
		return errors.Wrap(err)
	}

	cart, err := h.cartService.GetCartById(msgIn.Ctx, cartId)
	if err != nil {
		return errors.Wrap(err)
	}

	msg := message.OpenCart(msgIn.Ctx, cart)
	err = out.SendMessage(msg)
	if err != nil {
		return errors.Wrap(err)
	}

	err = h.cartService.SyncCartMessage(msgIn.Ctx, cart, msg)
	if err != nil {
		return errors.Wrap(err)
	}

	if !msgIn.IsCallback {
		_ = out.SendMessage(&response.DeleteMessage{
			ChatId:    msgIn.Chat.ID,
			MessageId: int64(msgIn.MessageID),
		})
	}

	return nil
}

// GetCommand - возвращает комманду
func (h *Handler) GetCommand() string {
	return commands.Purge
}
