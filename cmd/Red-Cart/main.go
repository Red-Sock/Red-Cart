package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/Red-Sock/Red-Cart/internal/clients/telegram"
	"github.com/Red-Sock/Red-Cart/internal/config"
	telegramserver "github.com/Red-Sock/Red-Cart/internal/transport/telegram"
	"github.com/Red-Sock/Red-Cart/internal/utils/closer"
)

func main() {
	logrus.Println("starting app")

	ctx := context.Background()

	cfg, err := config.ReadConfig()
	if err != nil {
		logrus.Fatalf("error reading config %s", err.Error())
	}

	startupDuration, err := cfg.GetDuration(config.AppInfoStartupDuration)
	if err != nil {
		logrus.Fatalf("error extracting startup duration %s", err)
	}
	ctx, cancel := context.WithTimeout(ctx, startupDuration)

	closer.Add(func() error {
		cancel()
		return nil
	})

	tg := telegramserver.NewServer(cfg, telegram.New(cfg))
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
