package commands

func NewDeleteCommand(itemKey string) string {
	return DeleteItem + " " + itemKey
}

func NewCheckCommand(itemKey string) string {
	return CheckItem + " " + itemKey
}

func NewUncheckCommand(itemKey string) string {
	return UncheckItem + " " + itemKey
}
