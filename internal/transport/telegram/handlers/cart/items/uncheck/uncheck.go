package uncheck

import (
	"strings"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
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

// Handle - expects to have cart id and item name as an input argument
func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) error {
	if len(in.Args) < 1 {
		return out.SendMessage(response.NewMessage("expect to have 2 argements as input - cart id and item name"))
	}

	userCart, err := h.cartService.GetCartByChatId(in.Ctx, in.Chat.ID)
	if err != nil {
		return errors.Wrap(err)
	}

	itemKey := strings.Join(in.Args, " ")

	err = h.itemService.Uncheck(in.Ctx, userCart.Cart.Id, itemKey)
	if err != nil {
		return errors.Wrap(err)
	}

	userCart, err = h.cartService.GetCartById(in.Ctx, userCart.Cart.Id)
	if err != nil {
		return errors.Wrap(err)
	}

	msg := message.OpenCart(in.Ctx, userCart)
	err = out.SendMessage(msg)
	if err != nil {
		return errors.Wrap(err)
	}

	err = h.cartService.SyncCartMessage(in.Ctx, userCart, msg)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (h *Handler) GetCommand() string {
	return commands.UncheckItem
}
