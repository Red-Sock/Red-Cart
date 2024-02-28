package tests

import (
	"context"
	_ "embed"
	"fmt"
	"sync/atomic"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/jackc/pgx/v5"
	"github.com/pressly/goose/v3"

	postgresclient "github.com/Red-Sock/Red-Cart/internal/clients/postgres"
	"github.com/Red-Sock/Red-Cart/internal/config"
	"github.com/Red-Sock/Red-Cart/internal/data/postgres"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/data"
	"github.com/Red-Sock/Red-Cart/internal/interfaces/service"
	servicev1 "github.com/Red-Sock/Red-Cart/internal/service"
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

func UseServiceV1(a *App) {
	a.Srv = servicev1.New(a.Db)
}

var pg data.Storage

func UsePgDb(a *App) {
	if pg != nil {
		a.Db = pg
		return
	}

	cfg := getConfig()
	ctx := context.Background()

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.GetString(config.DataSourcesPostgresUser),
		cfg.GetString(config.DataSourcesPostgresPwd),
		cfg.GetString(config.DataSourcesPostgresHost),
		cfg.GetString(config.DataSourcesPostgresPort),
		cfg.GetString(config.DataSourcesPostgresName),
	)

	pgxConn, err := pgx.Connect(ctx, connString)
	if err != nil {
		panic(err)
	}

	_, err = pgxConn.Exec(ctx, `
	DROP SCHEMA public CASCADE;
	CREATE SCHEMA public;
	`)
	if err != nil {
		panic(err)
	}

	conn, err := goose.OpenDBWithDriver("pgx", connString)
	if err != nil {
		panic(errors.Wrap(err, "error opening pg db connection"))
	}

	err = goose.Up(conn, "./../migrations")
	if err != nil {
		panic(err)
	}

	txMgr := postgresclient.NewTx(pgxConn)
	pg = postgres.New(txMgr)

	a.Db = pg
}

var (
	//go:embed tests_config.yaml
	cfgFile []byte

	conf *config.Config
)

func getConfig() *config.Config {
	if conf != nil {
		return conf
	}

	var err error
	conf, err = config.ParseConfig(cfgFile)
	if err != nil {
		panic(errors.Wrap(err, "error parsing test config"))
	}

	return conf
}

var userIdGenerator atomic.Int64

func GetUserID() int64 {
	return userIdGenerator.Add(1)
}
