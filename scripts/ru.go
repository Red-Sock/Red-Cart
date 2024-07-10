package scripts

var ru = map[PhraseKey]string{
	OpenSetting: "Настроить корзину " + SettingIcon,
	CreateCart:  "Создать корзину " + CartIcon,
	Clear:       "Удалить товар / очистить корзину " + BinIcon,
	PurgeCart:   "Очистить корзину" + PurgeIcon,
	CartIsEmpty: "Корзина пуста",

	ClickToRemove: "Нажмите для удаления",
	Cart:          "Корзина" + CartIcon,

	Rename: "Переименовать " + EditIcon,

	Welcome:     "Добро пожаловать!",
	WelcomeBack: "С возвращением!",

	WelcomeMessagePattern: CartIcon + `
Корзина по умолчанию: %d

Для добавления продуктов просто введите их название
`,
}
