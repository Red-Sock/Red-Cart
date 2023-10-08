package welcome

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
)

const Command = "/welcome"

type Handler struct {
}

func (h *Handler) GetDescription() string {
	return "returns just hello"
}

func (h *Handler) GetCommand() string {
	return Command
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {
	out.SendMessage(response.NewMessage("Hello user!"))
}
