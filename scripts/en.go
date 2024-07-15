package scripts

import (
	"github.com/Red-Sock/Red-Cart/scripts/icons"
)

var (
	en = map[PhraseKey]string{
		CreateCartAction: "Create cart " + icons.CartIcon,
		OpenClearMenu:    "Remove items / purge cart " + icons.BinIcon,
		PurgeCartAction:  "Purge cart" + icons.PurgeIcon,
		CartIsEmpty:      "Cart is empty. Just in product names in order to add it to cart",

		Cart: "Cart " + icons.CartIcon,

		Rename: "Rename  " + icons.EditIcon,

		Welcome:     "Welcome!",
		WelcomeBack: "Welcome back!",

		WelcomeMessagePattern: icons.CartIcon + `
Default cart: %d
`,
	}
)
