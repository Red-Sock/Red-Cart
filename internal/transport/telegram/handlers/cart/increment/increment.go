package increment

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"

	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
)

const Command = "/add"

type Handler struct {
	service.CartService
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {

}

func (h *Handler) GetDescription() string {
	return "Increments something on given amount (or 1 by default"
}

func (h *Handler) GetCommand() string {
	return Command
}
