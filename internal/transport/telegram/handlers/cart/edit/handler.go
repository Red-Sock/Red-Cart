package edit

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"

	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
)

const Command = "/edit"

type Handler struct {
	cartService service.CartService
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {

}

func (h *Handler) GetDescription() string {
	return "Edit item in cart"
}

func (h *Handler) GetCommand() string {
	return Command
}
