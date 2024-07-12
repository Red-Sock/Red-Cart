package default_handler

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/helpers"
	"github.com/Red-Sock/Red-Cart/scripts"
)

type DefaultHandler struct {
	userService service.UserService
	cartService service.CartService
	itemService service.ItemService

	// lang ->> instruction name on native language ->> instruction
	expectedInstructions map[string]map[string]scripts.PhraseKey

	handlers map[string]tgapi.CommandHandler
}

func NewDefaultCommandHandler(
	srv service.Service,
	handlers map[string]tgapi.CommandHandler,
) *DefaultHandler {
	return &DefaultHandler{
		userService: srv.User(),
		cartService: srv.Cart(),
		itemService: srv.Item(),

		expectedInstructions: scripts.GetInstructions(),

		handlers: handlers,
	}
}

func (d *DefaultHandler) Handle(msgIn *model.MessageIn, out tgapi.Chat) error {
	defer helpers.DeleteIncomingMessage(msgIn, out)

	if len(msgIn.Args) == 0 || msgIn.Command != "" {
		err := out.SendMessage(response.NewMessage("unknown functionality " + msgIn.Command))
		if err != nil {
			return errors.Wrap(err)
		}

		return nil
	}

	handler, err := d.pickHandler(msgIn)
	if err != nil {
		return errors.Wrap(err)
	}

	if handler == nil {
		return nil
	}

	err = handler.Handle(msgIn, out)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (d *DefaultHandler) cartInputs(in *model.MessageIn, userCart domain.UserCart) (tgapi.MessageOut, error) {
	switch userCart.Cart.State {
	case domain.CartStateAdding:
		return d.addItem(in, userCart)
	case domain.CartStateEditingItemName:
		return d.editItemName(in, userCart)
	default:
		return nil, nil
	}
}
