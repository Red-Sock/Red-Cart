package delete_item

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

// Handle expects cart id and item name as arguments
func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {
	if len(in.Args) < 2 {
		out.SendMessage(response.NewMessage("expects cart id and item name as arguments"))
		return
	}

	cartId, err := strconv.ParseInt(in.Args[0], 10, 64)
	if err != nil {
		out.SendMessage(response.NewMessage("cart id must be integer."))
		return
	}

	err = h.itemService.Delete(in.Ctx, cartId, in.Args[1])
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	cart, err := h.cartService.GetCartById(in.Ctx, cartId)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	if !in.IsCallback {
		out.SendMessage(&response.DeleteMessage{
			ChatId:    in.Chat.ID,
			MessageId: *cart.Cart.MessageID,
		})
		message.CartFromDomain(in.Ctx, out, cart)
		return
	}

	if len(cart.Cart.Items) != 0 {
		message.Delete(out, cart)
		return
	}

	message.CartFromDomain(in.Ctx, out, cart)
}

func (h *Handler) GetCommand() string {
	return commands.DeleteItem
}
