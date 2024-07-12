package main

import (
	"github.com/sirupsen/logrus"

	"github.com/Red-Sock/Red-Cart/cmd/Red-Cart/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		logrus.Fatal(err)
		return
	}

	err = a.Start()
	if err != nil {
		logrus.Fatal(err)
	}
}
