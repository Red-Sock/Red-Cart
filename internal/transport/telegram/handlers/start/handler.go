package start

import (
	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
)

const Command = "/start"

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
	newUser := domain.User{
		Id:        in.From.ID,
		UserName:  in.From.UserName,
		FirstName: in.From.FirstName,
		LastName:  in.From.LastName,
	}
	startMessage, err := h.userSrv.Start(in.Ctx, newUser)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	msg := response.NewMessage(startMessage)
	msg.Keys = &keyboard.InlineKeyboard{}

	msg.Keys.AddButton("Создать корзину", "/create_cart")
	msg.Keys.AddButton("Добавить товар", "/add_item")
	out.SendMessage(msg)
}
