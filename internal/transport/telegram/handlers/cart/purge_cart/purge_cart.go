package purge_cart

import (
	"strconv"

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

// Handle expects to have cart id as a given parameter
func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) error {
	if len(in.Args) < 1 {
		return out.SendMessage(response.NewMessage("expecting to have a cart id as an argument"))
	}

	cartId, err := strconv.ParseInt(in.Args[0], 10, 64)
	if err != nil {
		return out.SendMessage(response.NewMessage("expected for cart id to be integer: " + err.Error()))
	}

	err = h.cartService.PurgeCart(in.Ctx, cartId)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	cart, err := h.cartService.GetCartById(in.Ctx, cartId)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	msg, err := message.OpenCart(in.Ctx, out, cart)
	if err != nil {
		return err
	}

	err = h.cartService.SyncCartMessage(in.Ctx, cart, msg)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	if !in.IsCallback {
		_ = out.SendMessage(&response.DeleteMessage{
			ChatId:    in.Chat.ID,
			MessageId: int64(in.MessageID),
		})
	}

	return nil
}

func (h *Handler) GetCommand() string {
	return commands.Purge
}
