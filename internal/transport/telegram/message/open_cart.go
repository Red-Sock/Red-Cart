package message

import (
	"context"
	"strconv"

	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/response"

	"github.com/Red-Sock/Red-Cart/internal/domain"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/scripts"
)

func OpenCart(ctx context.Context, userCart domain.UserCart) interfaces.MessageOut {
	var text string

	if len(userCart.Cart.Items) == 0 {
		return emptyCart(ctx, userCart)
	}

	text = scripts.Get(ctx, scripts.Cart)
	keys := CartKeys(userCart.Cart)

	if userCart.Cart.MessageId != nil {
		return &response.EditMessage{
			ChatId:    userCart.Cart.ChatId,
			MessageId: *userCart.Cart.MessageId,
			Text:      text,
			Keys:      keys,
		}
	}

	return &response.MessageOut{
		ChatId: userCart.User.Id,
		Text:   text,
		Keys:   keys,
	}
}

func CartSettings(ctx context.Context, cart domain.UserCart) interfaces.MessageOut {
	var text string
	if len(cart.Cart.Items) == 0 {
		text = scripts.Get(ctx, scripts.CartIsEmpty)
	} else {
		text = scripts.Get(ctx, scripts.Cart)
	}

	var keys *keyboard.GridKeyboard

	if len(cart.Cart.Items) != 0 {
		keys = &keyboard.GridKeyboard{}
		keys.Columns = 1

		itemsNames, itemKeys := itemList(cart.Cart.Items)
		for i, itemName := range itemsNames {
			keys.AddButton(itemName, commands.EditItem+" "+itemKeys[i])
		}

		keys.AddButton("️🔙", commands.Cart)
	}

	if cart.Cart.MessageId != nil {
		return &response.EditMessage{
			ChatId:    cart.Cart.ChatId,
			MessageId: *cart.Cart.MessageId,
			Text:      text,
			Keys:      keys,
		}
	}

	out := &response.MessageOut{
		ChatId: cart.User.Id,
		Text:   text,
		Keys:   keys,
	}

	return out
}

func emptyCart(ctx context.Context, cart domain.UserCart) interfaces.MessageOut {
	text := scripts.Get(ctx, scripts.CartIsEmpty)

	var msg interfaces.MessageOut
	if cart.Cart.MessageId != nil {
		msg = &response.EditMessage{
			ChatId:    cart.Cart.ChatId,
			Text:      text,
			MessageId: *cart.Cart.MessageId,
		}
	} else {
		msg = response.NewMessage(text)
	}

	return msg
}

func CartKeys(cart domain.Cart) (keys *keyboard.FloatingKeyboard) {
	cartId := strconv.Itoa(int(cart.ID))
	if len(cart.Items) == 0 {
		return nil
	}
	keys = &keyboard.FloatingKeyboard{}

	items, itemKeys := itemList(cart.Items)
	for i := range items {
		row := make([]keyboard.Button, 0, 4)

		row = append(row, keyboard.Button{
			Text:  "+",
			Value: "+", // todo
		})

		row = append(row, getItemButton(cart.Items[i], cartId, itemKeys[i]))

		row = append(row, keyboard.Button{
			Text:  "-",
			Value: "-", // todo
		})

		row = append(row, keyboard.Button{
			Text:  scripts.BinIcon,
			Value: "delete", // todo
		})
		keys.AddRow(row)
	}

	return keys
}

func getItemButton(item domain.Item, cartId, itemKey string) (key keyboard.Button) {
	key.Value = commands.CheckItem

	key.Text = item.Name
	if item.Amount > 1 {
		key.Text += "(" + strconv.Itoa(int(item.Amount)) + ")"
	}

	if item.Checked {
		key.Text += scripts.CheckedIcon
		key.Value = commands.UncheckItem
	}

	key.Value += " " + itemKey

	return key
}

func getAddItemButton(item domain.Item) {

}
