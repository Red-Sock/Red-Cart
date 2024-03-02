package handlers

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/scripts"
)

type DefaultHandler struct {
	userService service.UserService
	cartService service.CartService
	itemService service.ItemService

	expectedInstructions map[string]map[string]scripts.PhraseKey // lang ->> instruction name on native language ->> instruction
}

func NewDefaultCommandHandler(
	us service.UserService,
	cs service.CartService,
	is service.ItemService,
) *DefaultHandler {
	return &DefaultHandler{
		userService:          us,
		cartService:          cs,
		itemService:          is,
		expectedInstructions: scripts.GetInstructions(),
	}
}

func (d *DefaultHandler) Handle(in *model.MessageIn, out tgapi.Chat) {
	if len(in.Args) == 0 || in.Command != "" {
		out.SendMessage(response.NewMessage("unknown functionality " + in.Command))
		return
	}

	if d.basicInputs(in, out) {
		return
	}

	if d.cartInputs(in, out) {
		return
	}
}

func (d *DefaultHandler) cartInputs(in *model.MessageIn, out tgapi.Chat) bool {
	userCart, err := d.cartService.GetCartByChatId(in.Ctx, in.Chat.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return true
	}

	switch userCart.Cart.State {
	case domain.CartStateAdding:
		d.addItem(in, out, userCart)
	case domain.CartStateEditingItemName:
		d.editItemName(in, out, userCart)
	}

	return false
}
