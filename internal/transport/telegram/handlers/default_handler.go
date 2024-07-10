package handlers

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/scripts"
)

type DefaultHandler struct {
	userService service.UserService
	cartService service.CartService
	itemService service.ItemService

	// lang ->> instruction name on native language ->> instruction
	expectedInstructions map[string]map[string]scripts.PhraseKey
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

func (d *DefaultHandler) Handle(msgIn *model.MessageIn, out tgapi.Chat) error {
	if len(msgIn.Args) == 0 || msgIn.Command != "" {
		err := out.SendMessage(response.NewMessage("unknown functionality " + msgIn.Command))
		if err != nil {
			return errors.Wrap(err)
		}

		return nil
	}

	userCart, err := d.cartService.GetCartByChatId(msgIn.Ctx, msgIn.Chat.ID)
	if err != nil {
		err = out.SendMessage(response.NewMessage(err.Error()))
		if err != nil {
			return errors.Wrap(err)
		}

		return nil
	}

	ok, err := d.basicInputs(msgIn, userCart, out)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	ok, err = d.cartInputs(msgIn, userCart, out)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	return nil
}

func (d *DefaultHandler) cartInputs(in *model.MessageIn, userCart domain.UserCart, out tgapi.Chat) (bool, error) {
	switch userCart.Cart.State {
	case domain.CartStateAdding:
		return true, d.addItem(in, out, userCart)
	case domain.CartStateEditingItemName:
		return true, d.editItemName(in, out, userCart)
	}

	return false, nil
}
