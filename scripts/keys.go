package scripts

const (
	CreateCartAction PhraseKey = "CreateCart"
	PurgeCartAction  PhraseKey = "Purge"

	OpenClearMenu PhraseKey = "Clear"

	Cart   PhraseKey = "Cart"
	Rename PhraseKey = "Rename"

	Welcome               PhraseKey = "Welcome"
	WelcomeBack           PhraseKey = "WelcomeBack"
	WelcomeMessagePattern PhraseKey = "WelcomeMessagePattern"
)

const (
	CartIsEmpty PhraseKey = "CartIsEmpty"
)

type PhraseKey string
