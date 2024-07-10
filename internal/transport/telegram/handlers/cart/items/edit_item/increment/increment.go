package increment

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"

	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) error {
	return nil
}

func (h *Handler) GetCommand() string {
	return commands.Increment
}
