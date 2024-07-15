package delete_item

import (
	"strings"

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

	if len(msgIn.Args) < 1 {
		return out.SendMessage(response.NewMessage("expects cart id and item name as arguments"))
	}

	itemName := strings.Join(msgIn.Args, " ")
	userCart, err := h.cartService.GetCartByChatId(msgIn.Ctx, msgIn.Chat.ID)
	if err != nil {
		return errors.Wrap(err)
	}

	err = h.itemService.Delete(msgIn.Ctx, userCart.Cart.Id, itemName)
	if err != nil {
		return errors.Wrap(err)
	}

	userCart, err = h.cartService.GetCartById(msgIn.Ctx, userCart.Cart.Id)
	if err != nil {
		return errors.Wrap(err)
	}

	msg := message.ClearCartMenu(msgIn.Ctx, userCart)
	err = out.SendMessage(msg)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (h *Handler) GetCommand() string {
	return commands.DeleteItem
}
