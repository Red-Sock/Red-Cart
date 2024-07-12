package telegram

import (
	"github.com/Red-Sock/go_tg/interfaces"

	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/commands"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/items/add"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/items/check"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/items/delete_item"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/items/edit_item"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/items/edit_item/increment"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/items/edit_item/rename"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/items/uncheck"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/open_clear_menu"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/purge_cart"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/cart/settings"
	"github.com/Red-Sock/Red-Cart/internal/transport/telegram/handlers/start"
)

type HandlerStore struct {
	handlers map[string]interfaces.CommandHandler
}

func newHandlerStore(srv service.Service) *HandlerStore {
	return &HandlerStore{
		handlers: map[string]interfaces.CommandHandler{
			commands.Start: start.New(srv),
			commands.Cart:  cart.New(srv),

			commands.AddItem:            add.New(srv),
			commands.EditItem:           edit_item.New(srv),
			commands.RenameItem:         rename.New(srv),
			commands.IncrementItemCount: increment.New(),
			commands.DeleteItem:         delete_item.New(srv),

			commands.CheckItem:   check.New(srv),
			commands.UncheckItem: uncheck.New(srv),

			commands.ClearMenu: open_clear_menu.New(srv),

			commands.CartSetting: settings.New(srv),

			commands.Purge: purge_cart.New(srv),
		},
	}
}
