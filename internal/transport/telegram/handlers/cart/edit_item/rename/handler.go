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
func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {
	if len(in.Args) < 2 {
		out.SendMessage(response.NewMessage("Expected to have cart id and item name as arguments"))
		return
	}

	cartID, err := strconv.ParseInt(in.Args[0], 10, 64)
	if err != nil {
		out.SendMessage(response.NewMessage("Expected to have cart id as integer type:" + err.Error()))
		return
	}

	cart, err := h.cartService.GetCartById(in.Ctx, cartID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	var itemInCart *domain.Item

	for idx := range cart.Cart.Items {
		if cart.Cart.Items[idx].Name == in.Args[1] {
			itemInCart = &cart.Cart.Items[idx]
		}
	}

	if itemInCart == nil {
		out.SendMessage(response.NewMessage("no item with such name"))
		return
	}

	msg := h.sendRenameMessage(cart.Cart, itemInCart.Name, out)

	err = h.cartService.SyncCartMessage(in.Ctx, cart.Cart, msg)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	err = h.cartService.AwaitNameChange(in.Ctx, cartID, *itemInCart)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}
}

func (h *Handler) sendRenameMessage(cart domain.Cart, oldItemName string, out tgapi.Chat) tgapi.MessageOut {
	text := "Введите новое имя для " + oldItemName

	if cart.MessageID != nil {
		msg := &response.EditMessage{
			ChatId:    cart.ChatID,
			Text:      text,
			MessageId: *cart.MessageID,
		}
		out.SendMessage(msg)

		if msg.GetMessageId() != 0 {
			return msg
		}
	}

	msg := response.NewMessage(text)

	return msg
}

func (h *Handler) GetCommand() string {
	return commands.Rename
}
