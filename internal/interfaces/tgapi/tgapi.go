package tgapi

import (
	gotg "github.com/Red-Sock/go_tg/interfaces"
)

type TgApi interface {
	Start() error
	Stop()
	AddCommandHandler(handler gotg.CommandHandler)
	SetDefaultCommandHandler(h gotg.CommandHandler)
}
