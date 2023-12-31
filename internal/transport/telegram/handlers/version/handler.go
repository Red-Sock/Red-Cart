package version

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/config"
)

const Command = "/version"

type Handler struct {
	version string
}

func New(cfg *config.Config) *Handler {
	return &Handler{
		version: cfg.GetString(config.AppInfoVersion),
	}
}

func (h *Handler) GetDescription() string {
	return "returns current app version as a response"
}

func (h *Handler) GetCommand() string {
	return Command
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {
	out.SendMessage(response.NewMessage(in.Text + ": " + h.version))
}
