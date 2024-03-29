package telegram

import (
	"github.com/Red-Sock/go_tg"
	"github.com/godverv/matreshka/resources"
)

func New(cfg *resources.Telegram) *go_tg.Bot {
	return go_tg.NewBot(cfg.ApiKey)
}
