package scripts

var en = map[PhraseKey]string{
	OpenSetting: "Cart settings " + SettingIcon,
	CreateCart:  "Create cart " + CartIcon,
	Clear:       "Remove items / purge cart " + BinIcon,
	PurgeCart:   "Purge cart" + PurgeIcon,
	CartIsEmpty: "Cart is empty",

	ClickToRemove: "Click to remove",
	Cart:          "Cart " + CartIcon,

	Rename: "Rename  " + EditIcon,

	Welcome:     "Welcome!",
	WelcomeBack: "Welcome back!",

	WelcomeMessagePattern: CartIcon + `
Default cart: %d

Just in product names in order to add it to cart
`,
}
