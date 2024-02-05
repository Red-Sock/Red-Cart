package checkout

import (
	"strconv"
	"strings"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
)

const Command = "/checkout"

type Handler struct {
	cartService service.CartService
}

func New(service service.CartService) *Handler {
	return &Handler{
		cartService: service,
	}
}

func (h *Handler) GetDescription() string {
	return "Show items in your  cart"
}

func (h *Handler) GetCommand() string {
	return Command
}

func (h *Handler) Handle(in *model.MessageIn, out tgapi.Chat) {
	_, err := h.cartService.GetByOwnerId(in.Ctx, in.From.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	cartItem, err := h.cartService.ShowCartItem(in.Ctx, in.From.ID)
	var outMessageBuilder strings.Builder
	for _, item := range cartItem {
		outMessageBuilder.WriteString("User: ")
		outMessageBuilder.WriteString(strconv.FormatInt(item.UserID, 10))
		outMessageBuilder.WriteString("\n")
		for _, name := range item.ItemNames {
			outMessageBuilder.WriteString(name)
			outMessageBuilder.WriteString("\n")
		}
	}

	out.SendMessage(response.NewMessage(outMessageBuilder.String()))

}
