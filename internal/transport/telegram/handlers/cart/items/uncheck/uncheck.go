package uncheck

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
	itemService service.ItemService
	cartService service.CartService
}

func New(itemService service.ItemService, cartService service.CartService) *Handler {
	return &Handler{
		itemService: itemService,
		cartService: cartService,
	}
}

// nolint
// Handle - expects to have cart id and item name as an input argument
func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) error {
	if len(in.Args) < 2 {
		return out.SendMessage(response.NewMessage("expect to have 2 argements as input - cart id and item name"))
	}

	cartId, err := strconv.ParseInt(in.Args[0], 10, 64)
	if err != nil {
		return out.SendMessage(response.NewMessage("cart id must be integer:" + err.Error()))
	}

	err = h.itemService.Uncheck(in.Ctx, cartId, in.Args[1])
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	userCart, err := h.cartService.GetCartById(in.Ctx, cartId)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	msg := message.OpenCart(in.Ctx, userCart)
	err = out.SendMessage(msg)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	err = h.cartService.SyncCartMessage(in.Ctx, userCart, msg)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	return nil
}

func (h *Handler) GetCommand() string {
	return commands.Uncheck
}
