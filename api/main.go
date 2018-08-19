package main

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/sirkuttin/acme-cert-api/api/api"
	"os"
)

var log = logrus.New()

func init() {
	log.SetLevel(logrus.DebugLevel)
}

func main() {

	errChan := make(chan error)

	go func() {
		api.Start(log)
		errChan <- errors.New("api exited")
	}()

	log.Error("error: ", <-errChan)
	os.Exit(1)
}
