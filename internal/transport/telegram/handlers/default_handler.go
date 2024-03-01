package handlers

import (
	"encoding/json"

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
	itemService service.ItemService
}

func NewDefaultCommandHandler(
	us service.UserService,
	cs service.CartService,
	is service.ItemService,
) *DefaultHandler {
	return &DefaultHandler{
		userService: us,
		cartService: cs,
		itemService: is,
	}
}

func (d *DefaultHandler) Handle(in *model.MessageIn, out tgapi.Chat) {
	if len(in.Args) == 0 || in.Command != "" {
		return
	}

	userCart, err := d.cartService.GetCartByChatId(in.Ctx, in.Chat.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	switch userCart.Cart.State {
	case domain.CartStateAdding:
		d.addItem(in, out, userCart)
	case domain.CartStateEditingItemName:
		d.editItemName(in, out, userCart)

	}
}

func (d *DefaultHandler) addItem(in *model.MessageIn, out tgapi.Chat, userCart domain.UserCart) {
	items := make([]domain.Item, 0, len(in.Args))
	for _, item := range in.Args {
		items = append(items, domain.Item{Name: item, Amount: 1})
	}

	cart, err := d.cartService.Add(in.Ctx, items, userCart.Cart.ID, in.From.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	out.SendMessage(&response.DeleteMessage{
		ChatId:    in.Chat.ID,
		MessageId: int64(in.MessageID),
	})

	msg := message.CartFromDomain(out, cart)

	err = d.cartService.SyncCartMessage(in.Ctx, userCart.Cart, msg)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}
}

func (d *DefaultHandler) editItemName(in *model.MessageIn, out tgapi.Chat, cart domain.UserCart) {
	if in.Text == "" {
		out.SendMessage(response.NewMessage("in order to change name you have to pass a valid string name"))
		return
	}
	var p domain.ChangeItemNamePayload
	err := json.Unmarshal(cart.Cart.StatePayload, &p)
	if err != nil {
		out.SendMessage(response.NewMessage("error parsing cart payload"))
		return
	}

	out.SendMessage(&response.DeleteMessage{
		ChatId:    in.Chat.ID,
		MessageId: int64(in.MessageID),
	})

	err = d.itemService.UpdateName(in.Ctx, cart.Cart.ID, p.ItemName, in.Text)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	err = d.cartService.AwaitItemsAdded(in.Ctx, cart.Cart.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	cart, err = d.cartService.GetCartById(in.Ctx, cart.Cart.ID)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}

	msg := message.CartFromDomain(out, cart)

	err = d.cartService.SyncCartMessage(in.Ctx, cart.Cart, msg)
	if err != nil {
		out.SendMessage(response.NewMessage(err.Error()))
		return
	}
}
