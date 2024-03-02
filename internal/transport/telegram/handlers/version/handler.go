package version

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/config"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
)

type Handler struct {
	version string
}

func New(cfg *config.Config) *Handler {
	return &Handler{
		version: cfg.GetString(config.AppInfoVersion),
	}
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) error {
	return out.SendMessage(response.NewMessage(in.Text + ": " + h.version))
}

func (h *Handler) GetCommand() string {
	return commands.Version
}
