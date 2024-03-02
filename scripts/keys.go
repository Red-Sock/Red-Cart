package scripts

type PhraseKey string

const (
	OpenSetting PhraseKey = "Setting"
	CreateCart  PhraseKey = "CreateCart"
	Clear       PhraseKey = "Clear"
	PurgeCart   PhraseKey = "Purge"

	ClickToRemove PhraseKey = "ClickToRemove"
	Cart          PhraseKey = "Cart"
	Rename        PhraseKey = "Rename"
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
