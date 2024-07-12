package scripts

const (
	CreateCart PhraseKey = "CreateCart"
	Clear      PhraseKey = "Clear"
	PurgeCart  PhraseKey = "Purge"

	ClickToRemove PhraseKey = "ClickToRemove"
	Cart          PhraseKey = "Cart"
	Rename        PhraseKey = "Rename"

	Welcome               PhraseKey = "Welcome"
	WelcomeBack           PhraseKey = "WelcomeBack"
	WelcomeMessagePattern PhraseKey = "WelcomeMessagePattern"
)

const (
	CartIsEmpty PhraseKey = "CartIsEmpty"
)

const (
	CartIcon    = "🛒"
	SettingIcon = "🛠️"
	DeleteIcon  = "❌"
	BinIcon     = "🗑️"
	PurgeIcon   = "🚮"
	CheckedIcon = "✅"
	BackIcon    = "🔙"
	EditIcon    = "✏️"
)

type PhraseKey string
