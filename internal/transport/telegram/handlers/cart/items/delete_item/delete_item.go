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
func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) error {
	if len(in.Args) < 2 {
		return out.SendMessage(response.NewMessage("expects cart id and item name as arguments"))
	}

	cartId, err := strconv.ParseInt(in.Args[0], 10, 64)
	if err != nil {
		return out.SendMessage(response.NewMessage("cart id must be integer."))
	}

	err = h.itemService.Delete(in.Ctx, cartId, in.Args[1])
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	cart, err := h.cartService.GetCartById(in.Ctx, cartId)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	if !in.IsCallback {
		_ = out.SendMessage(&response.DeleteMessage{
			ChatId:    in.Chat.ID,
			MessageId: *cart.Cart.MessageID,
		})

		_, err = message.OpenCart(in.Ctx, out, cart)
		return err
	}

	if len(cart.Cart.Items) != 0 {
		_, err := message.Delete(out, cart)
		return err
	}

	_, err = message.OpenCart(in.Ctx, out, cart)
	return err
}

func (h *Handler) GetCommand() string {
	return commands.DeleteItem
}
