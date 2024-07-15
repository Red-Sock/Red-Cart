package scripts

import (
	"github.com/Red-Sock/Red-Cart/scripts/icons"
)

var (
	ru = map[PhraseKey]string{
		CreateCartAction: "Создать корзину " + icons.CartIcon,
		OpenClearMenu:    "Удалить товар / очистить корзину " + icons.BinIcon,
		PurgeCartAction:  "Очистить корзину" + icons.PurgeIcon,
		CartIsEmpty:      "Корзина пуста. Для добавления продуктов просто введите их название",

		Cart: "Корзина" + icons.CartIcon,

		Rename: "Переименовать " + icons.EditIcon,

		Welcome:     "Добро пожаловать!",
		WelcomeBack: "С возвращением!",

		WelcomeMessagePattern: icons.CartIcon + `
Корзина по умолчанию: %d
`,
	}
)
