package handlers

import (
	"strings"

	tgapi "github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/message"
)

type DefaultHandler struct {
	userService service.UserService
	cartService service.CartService
}

func NewDefaultCommandHandler(us service.UserService, cs service.CartService) *DefaultHandler {
	return &DefaultHandler{
		userService: us,
		cartService: cs,
	}
}

func (d *DefaultHandler) Handle(in *model.MessageIn, out tgapi.Chat) {
	argsIn := strings.Split(in.Text, "\n")
	items := make([]domain.Item, 0, len(argsIn))
	for _, item := range argsIn {
		items = append(items, domain.Item{Name: item, Amount: 1})
	}

	userCart, err := d.userService.AddToDefaultCart(in.Ctx, items, in.From.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	out.SendMessage(&response.DeleteMessage{
		ChatId:    in.Chat.ID,
		MessageId: int64(in.MessageID),
	})

	msg := message.CartFromDomain(userCart)

	oldChatID, oldMsgID := msg.GetChatId(), msg.GetMessageId()
	out.SendMessage(msg)
	newChatId, newMsgId := msg.GetChatId(), msg.GetMessageId()

	if oldChatID == newChatId && oldMsgID == newMsgId {
		return
	}

	userCart.Cart.ChatID = &newChatId
	userCart.Cart.MessageID = &newMsgId

	err = d.cartService.UpdateMessageRef(in.Ctx, userCart.Cart)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}
}

func (d *DefaultHandler) GetDescription() string {
	// TODO
	return "TODO"
}

func (d *DefaultHandler) GetCommand() string {
	// TODO
	return "TODO"
}
