package edit_item

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

type Handler struct {
	userService service.UserService
	cartService service.CartService
}

func New(userService service.UserService, cartService service.CartService) *Handler {
	return &Handler{
		userService: userService,
		cartService: cartService,
	}
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) error {
	if len(in.Args) == 0 {
		return nil
	}

	cart, err := h.userService.GetCartByChat(in.Ctx, in.Chat.ID)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	var itemInCart *domain.Item

	for idx := range cart.Cart.Items {
		if cart.Cart.Items[idx].Name == in.Args[0] {
			itemInCart = &cart.Cart.Items[idx]
		}
	}

	if itemInCart == nil {
		return nil
	}

	msg, err := message.EditFromCartItem(in.Ctx, out, cart, *itemInCart)
	if err != nil {
		return err
	}

	err = h.cartService.SyncCartMessage(in.Ctx, cart, msg)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	return nil
}

func (h *Handler) GetCommand() string {
	return commands.Edit
}
