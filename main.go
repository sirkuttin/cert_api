package main

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/sirkuttin/acme-cert-api/api"
	"os"
	"github.com/dgraph-io/badger"
	"syscall"
	"os/signal"
)

var log = logrus.New()

func init() {
	log.SetLevel(logrus.DebugLevel)
}

func main() {

	opts := badger.DefaultOptions
	opts.Dir = "app_db"
	opts.ValueDir = "app_db"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	errChan := make(chan error)
	go func() {
		api.Start(log, db)
		errChan <- errors.New("api exited")
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	var exitCode int
	select {
	case <-stop:
		log.Info("Graceful Shutdown")
		exitCode = 0
	case msg := <-errChan:
		log.Error("error: ", msg)
		exitCode = 1
	}

	db.Close()
	os.Exit(exitCode)
}
