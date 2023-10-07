package welcome

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/config"
)

const Command = "/welcome"

type Handler struct {
	version string
}

func (h *Handler) GetDescription() string {
	return "returns just hello"
}

func (h *Handler) GetCommand() string {
	return Command
}

func New(cfg *config.Config) *Handler {
	return &Handler{
		version: cfg.GetString(config.AppInfoVersion),
	}
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {
	out.SendMessage(response.NewMessage("Hello user!"))
}
