package add

import (
	"strings"

	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

type Handler struct {
	cartService service.CartService
}

func New(srv service.Service) *Handler {
	return &Handler{
		cartService: srv.Cart(),
	}
}

func (h *Handler) Handle(msgIn *model.MessageIn, out interfaces.Chat) error {
	itemsRaw := strings.Split(msgIn.Text, "\n")
	items := make([]domain.Item, 0, len(itemsRaw))

	for _, item := range itemsRaw {
		items = append(items, domain.Item{Name: item, Amount: 1})
	}

	userCart, err := h.cartService.GetCartByChatId(msgIn.Ctx, msgIn.Chat.ID)
	if err != nil {
		return errors.Wrap(err)
	}

	userCart, err = h.cartService.Add(msgIn.Ctx, items, userCart.Cart.ID, msgIn.From.ID)
	if err != nil {
		return errors.Wrap(err)
	}

	cartMessage := message.OpenCart(msgIn.Ctx, userCart)
	err = out.SendMessage(cartMessage)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (h *Handler) GetCommand() string {
	return commands.AddItem
}
