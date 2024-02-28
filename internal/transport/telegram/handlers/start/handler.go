package start

import (
	"fmt"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
)

const Command = "/start"

type Handler struct {
	userSrv service.UserService
	cartSrv service.CartService
}

func New(userSrv service.UserService, cartSrv service.CartService) *Handler {
	return &Handler{
		userSrv: userSrv,
		cartSrv: cartSrv,
	}
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {
	newUser := domain.User{
		ID:        in.From.ID,
		UserName:  in.From.UserName,
		FirstName: in.From.FirstName,
		LastName:  in.From.LastName,
	}
	startMessage, err := h.userSrv.Start(in.Ctx, newUser)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	msg := response.NewMessage(startMessage.Msg + fmt.Sprintf(` 🛒

Корзина по умолчанию: %d

Для добавления продуктов просто введите их название
`, startMessage.Cart.ID))

	out.SendMessage(msg)

	out.SendMessage(&response.DeleteMessage{
		ChatId:    in.Chat.ID,
		MessageId: int64(in.MessageID),
	})

	//out.SendMessage(cartMsg)

}

func (h *Handler) GetDescription() string {
	return "returns just hello"
}

func (h *Handler) GetCommand() string {
	return Command
}
