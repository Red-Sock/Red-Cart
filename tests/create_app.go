package tests

import (
	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
)

type appOption func(a *App)

type App struct {
	Db  data.Storage
	Srv service.Service
}

func CreateTestApp(options ...appOption) *App {
	a := &App{}

	for _, opt := range options {
		opt(a)
	}

	return a
}
