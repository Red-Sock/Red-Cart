package rename

import (
	"strconv"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
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

// Handle expects to have 2 arguments:
// in.Args[0] = cart id
// in.Args[1] = name of product in cart
func (h *Handler) Handle(msgIn *model.MessageIn, out tgapi.Chat) error {
	if len(msgIn.Args) < 2 {
		return out.SendMessage(response.NewMessage("Expected to have cart id and item name as arguments"))
	}

	cartID, err := strconv.ParseInt(msgIn.Args[0], 10, 64)
	if err != nil {
		return out.SendMessage(response.NewMessage("Expected to have cart id as integer type:" + err.Error()))
	}

	cart, err := h.cartService.GetCartById(msgIn.Ctx, cartID)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	var itemInCart *domain.Item

	for idx := range cart.Cart.Items {
		if cart.Cart.Items[idx].Name == msgIn.Args[1] {
			itemInCart = &cart.Cart.Items[idx]
		}
	}

	if itemInCart == nil {
		return out.SendMessage(response.NewMessage("no item with such name"))
	}

	msg, err := h.sendRenameMessage(cart.Cart, itemInCart.Name, out)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	err = h.cartService.SyncCartMessage(msgIn.Ctx, cart, msg)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	err = h.cartService.AwaitNameChange(msgIn.Ctx, cartID, *itemInCart)
	if err != nil {
		return out.SendMessage(response.NewMessage(err.Error()))
	}

	return nil
}

func (h *Handler) sendRenameMessage(cart domain.Cart, oldItemName string, out tgapi.Chat) (tgapi.MessageOut, error) {
	text := "Введите новое имя для " + oldItemName

	if cart.MessageId != nil {
		msg := &response.EditMessage{
			ChatId:    cart.ChatId,
			Text:      text,
			MessageId: *cart.MessageId,
		}
		err := out.SendMessage(msg)

		if err == nil {
			return msg, nil
		}
	}

	msg := response.NewMessage(text)

	return msg, out.SendMessage(msg)
}

func (h *Handler) GetCommand() string {
	return commands.Rename
}
