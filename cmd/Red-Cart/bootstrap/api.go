package bootstrap

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/Red-Sock/Red-Cart/internal/config"
	"github.com/Red-Sock/Red-Cart/internal/transport"
)

func ApiEntryPoint(ctx context.Context, cfg *config.Config) (func(context.Context) error, error) {
	mngr := transport.NewManager()

	go func() {
		err := mngr.Start(ctx)
		if err != nil {
			logrus.Fatalf("error starting server %s", err.Error())
		}
	}()

	return mngr.Stop, nil
}
