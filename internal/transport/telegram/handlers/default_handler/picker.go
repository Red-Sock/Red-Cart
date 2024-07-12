package default_handler

import (
	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/scripts"
)

func (d *DefaultHandler) pickHandler(msgIn *model.MessageIn) (interfaces.Handler, error) {
	basicHandler := d.pickBasicHandler(msgIn)
	if basicHandler != nil {
		return basicHandler, nil
	}

	userCart, err := d.cartService.GetCartByChatId(msgIn.Ctx, msgIn.Chat.ID)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	cartHandler := d.pickCartHandler(userCart.Cart)
	if cartHandler != nil {
		return cartHandler, nil
	}

	return nil, nil
}

func (d *DefaultHandler) pickBasicHandler(msgIn *model.MessageIn) interfaces.Handler {
	lang := scripts.GetLangFromCtx(msgIn.Ctx)
	langInstr, ok := d.expectedInstructions[string(lang)]
	if !ok {
		return nil
	}

	instruction, ok := langInstr[msgIn.Text]
	if !ok {
		return nil
	}

	switch instruction {
	case scripts.OpenClearMenu:
		return d.handlers[commands.ClearMenu]
	default:
		return nil
	}
}

func (d *DefaultHandler) pickCartHandler(cart domain.Cart) interfaces.Handler {
	switch cart.State {
	case domain.CartStateAdding:
		return d.handlers[commands.AddItem]
	default:
		return nil
	}
}
