package decline

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
)

const Command = "/decline"

type Handler struct {
	cartService service.CartService
}

func New(service service.CartService) *Handler {
	return &Handler{
		cartService: service,
	}
}

func (h *Handler) GetDescription() string {
	return "decline order"
}

func (h *Handler) GetCommand() string {
	return Command
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {

	out.SendMessage(response.NewMessage("Ну как хочешь("))

}
