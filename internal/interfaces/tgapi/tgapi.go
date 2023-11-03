package tgapi

import (
	gotg "github.com/Red-Sock/go_tg/interfaces"
)

type TgApi interface {
	Start()
	Stop()
	AddCommandHandler(handler gotg.CommandHandler)
}
