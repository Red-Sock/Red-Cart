package checkout

import (
	"strconv"
	"strings"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/keyboard"
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

	usersItems, err := h.cartService.ShowCartItem(in.Ctx, in.From.ID)
	var outMessageBuilder strings.Builder
	for userID, items := range usersItems {
		outMessageBuilder.WriteString("User: ")
		outMessageBuilder.WriteString(strconv.FormatInt(userID, 10))
		outMessageBuilder.WriteString(" 🕐\n")

		for _, item := range items {
			outMessageBuilder.WriteString(item.Name)
			outMessageBuilder.WriteString("\n")
		}

		msg := response.NewMessageToChat("Подтверждать будешь пес?", userID)
		msg.Keys = &keyboard.InlineKeyboard{}

		msg.Keys.AddButton("✅", "/accept")
		msg.Keys.AddButton("❌", "/decline")
		out.SendMessage(msg)
	}

	out.SendMessage(response.NewMessage(outMessageBuilder.String()))

}
