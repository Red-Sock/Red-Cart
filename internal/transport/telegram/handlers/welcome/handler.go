package welcome

import (
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
)

const Command = "/welcome"

type Handler struct {
	userSrv service.UserService
}

func New(userSrv service.UserService) *Handler {
	return &Handler{
		userSrv: userSrv,
	}
}

func (h *Handler) GetDescription() string {
	return "returns just hello"
}

func (h *Handler) GetCommand() string {
	return Command
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {

	out.SendMessage(response.NewMessage(h.userSrv.Start(in.From.ID)))
}
