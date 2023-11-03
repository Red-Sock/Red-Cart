package tests

import (
	"github.com/Red-Sock/Red-Cart/internal/data/inmemory"
	"github.com/Red-Sock/Red-Cart/internal/service"
)

func UseInMemoryDb(a *App) {
	a.Db = inmemory.New()
}

func UseServiceV1(a *App) {
	a.Srv = service.New(a.Db)
}
