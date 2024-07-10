package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	pgclient "github.com/Red-Sock/Red-Cart/internal/clients/postgres"
	"github.com/Red-Sock/Red-Cart/internal/clients/telegram"
	"github.com/Red-Sock/Red-Cart/internal/config"
	"github.com/Red-Sock/Red-Cart/internal/data/postgres"
	"github.com/Red-Sock/Red-Cart/internal/service"
	telegramserver "github.com/Red-Sock/Red-Cart/internal/transport/telegram"
	"github.com/Red-Sock/Red-Cart/internal/utils/closer"
)

// nolint
func main() {
	logrus.Println("starting app")

	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("error reading config %s", err.Error())
	}

	startupDuration := cfg.GetAppInfo().StartupDuration
	if startupDuration == 0 {
		logrus.Fatalf("error extracting startup duration %s", err)
	}

	ctx, cancel := context.WithTimeout(ctx, startupDuration)

	closer.Add(func() error {
		cancel()

		return nil
	})

	p, err := cfg.GetDataSources().Postgres(config.ResourcePostgres)
	if err != nil {
		logrus.Fatalf("error getting postgres configuration %s", err)
	}

	conn, err := pgclient.New(ctx, p)
	if err != nil {
		logrus.Fatal(err, "error creating pgclient")
	}

	tgConf, err := cfg.GetDataSources().Telegram(config.ResourceTelegram)
	if err != nil {
		logrus.Fatal(err, "error getting telegram config")
	}

	tg := telegramserver.NewServer(
		cfg,
		telegram.New(tgConf),
		*service.New(postgres.New(conn)),
	)

	err = tg.Start(ctx)
	if err != nil {
		logrus.Fatalf("error starting telegram server %s", err)
	}

	waitingForTheEnd()

	err = tg.Stop(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Println("shutting down the app")

	if err = closer.Close(); err != nil {
		logrus.Fatalf("errors while shutting down application %s", err.Error())
	}
}

// rscli comment: an obligatory function for tool to work properly.
// must be called in the main function above
// also this is a LP song name reference, so no rules can be applied to the function name
func waitingForTheEnd() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
