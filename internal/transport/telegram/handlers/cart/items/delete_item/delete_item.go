package delete_item

import (
	"strconv"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/helpers"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

const (
	minArgumentsLength = 2

	cartIdIndex   = 0
	itemNameIndex = 1
)

type Handler struct {
	itemService service.ItemService
	cartService service.CartService
}

func New(srv service.Service) *Handler {
	return &Handler{
		itemService: srv.Item(),
		cartService: srv.Cart(),
	}
}

// Handle expects cart id and item name as arguments
func (h *Handler) Handle(msgIn *model.MessageIn, out tgapi.Chat) error {
	defer helpers.DeleteIncomingMessage(msgIn, out)

	if len(msgIn.Args) < minArgumentsLength {
		return out.SendMessage(response.NewMessage("expects cart id and item name as arguments"))
	}

	cartId, err := strconv.ParseInt(msgIn.Args[cartIdIndex], 10, 64)
	if err != nil {
		return out.SendMessage(response.NewMessage("cart id must be integer."))
	}

	err = h.itemService.Delete(msgIn.Ctx, cartId, msgIn.Args[itemNameIndex])
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	cart, err := h.cartService.GetCartById(msgIn.Ctx, cartId)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	msg := message.OpenCart(msgIn.Ctx, cart)
	err = out.SendMessage(msg)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (h *Handler) GetCommand() string {
	return commands.DeleteItem
}
