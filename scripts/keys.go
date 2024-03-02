package scripts

type PhraseKey string

const (
	OpenSetting PhraseKey = "Setting"
	CreateCart  PhraseKey = "CreateCart"
	Clear       PhraseKey = "Clear"
)

const (
	CartIsEmpty PhraseKey = "CartIsEmpty"
)

const (
	CartIcon    = "🛒"
	SettingIcon = "🛠️"
	DeleteIcon  = "❌"
	BinIcon     = "🗑️"
)
