package checkout

import (
	"strconv"
	"strings"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain/cart"
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
	cartOwner, err := h.cartService.GetByOwnerId(in.Ctx, in.From.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	cartItem, err := h.cartService.ShowCartItem(in.Ctx, in.From.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}
	var outMessageBuilder strings.Builder
	for _, item := range cartItem {

		outMessageBuilder.WriteString("User: ")
		outMessageBuilder.WriteString(strconv.FormatInt(item.UserID, 10))
		outMessageBuilder.WriteString(" 🕐\n")
		//Создаются товары конкретного пользователя
		var userItems strings.Builder
		for _, name := range item.ItemNames {
			userItems.WriteString(name)
			userItems.WriteString("\n")
		}
		outMessageBuilder.WriteString(userItems.String())
		//Создаем сообщение для подтверждения
		var newMessage strings.Builder
		newMessage.WriteString("User: ")
		newMessage.WriteString(strconv.FormatInt(cartOwner.OwnerId, 10))
		newMessage.WriteString(" просит подтвердить заказ:\n")
		newMessage.WriteString(userItems.String())
		msg := response.NewMessageToChat(newMessage.String(), item.UserID)
		msg.Keys = &keyboard.InlineKeyboard{}

		msg.Keys.AddButton("✅", "/accept")
		msg.Keys.AddButton("❌", "/decline")
		//Отправляем пользователю подтверждение
		out.SendMessage(msg)
	}
	//Отправляем владельцу
	out.SendMessage(response.NewMessage(outMessageBuilder.String()))

}

func checkStatus(cartItem []cart.CartItem) bool {
	for _, item := range cartItem {
		if item.Status == "" || item.Status == "wait" {
			return false
		}
	}

	return true
}
