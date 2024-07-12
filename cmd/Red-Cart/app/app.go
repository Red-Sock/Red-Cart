package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/sirupsen/logrus"

	pgclient "github.com/Red-Sock/Red-Cart/internal/clients/postgres"
	"github.com/Red-Sock/Red-Cart/internal/clients/telegram"
	"github.com/Red-Sock/Red-Cart/internal/config"
	"github.com/Red-Sock/Red-Cart/internal/data/postgres"
	"github.com/Red-Sock/Red-Cart/internal/service"
	telegramserver "github.com/Red-Sock/Red-Cart/internal/transport/telegram"
	"github.com/Red-Sock/Red-Cart/internal/utils/closer"
)

type App struct {
	server *telegramserver.Server
}

func New() (*App, error) {
	logrus.Println("starting app")

	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return nil, errors.Wrap(err, "error reading configs")
	}

	startupDuration := cfg.GetAppInfo().StartupDuration
	if startupDuration == 0 {
		return nil, errors.Wrap(err, "error extracting startup duration")
	}

	ctx, cancel := context.WithTimeout(ctx, startupDuration)

	closer.Add(func() error { cancel(); return nil })

	pg, err := cfg.GetDataSources().Postgres(config.ResourcePostgres)
	if err != nil {
		return nil, errors.Wrap(err, "error getting postgres configuration")
	}

	conn, err := pgclient.New(ctx, pg)
	if err != nil {
		return nil, errors.Wrap(err, "error creating sql connection")
	}

	tgConf, err := cfg.GetDataSources().Telegram(config.ResourceTelegram)
	if err != nil {
		logrus.Fatal(err, "error getting telegram config")
	}

	var app App
	app.server = telegramserver.NewServer(
		cfg,
		telegram.New(tgConf),
		service.New(postgres.New(conn)),
	)

	return &app, nil
}
func (a *App) Start() error {
	ctx := context.Background()
	err := a.server.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "error starting telegram server")
	}

	waitingForTheEnd()

	err = a.server.Stop(ctx)
	if err != nil {
		return errors.Wrap(err, "error stopping telegram server")
	}

	logrus.Println("shutting down the app")

	err = closer.Close()
	if err != nil {
		return errors.Wrap(err, "errors while shutting down application")
	}

	return nil
}

// rscli comment: an obligatory function for tool to work properly.
// must be called in the main function above
// also this is a LP song name reference, so no rules can be applied to the function name
func waitingForTheEnd() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
