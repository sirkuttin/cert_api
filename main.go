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
	opts.Dir = "tmp_db"
	opts.ValueDir = "tmp_db"
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

	select {
	case <-stop:
		log.Info("Graceful Shutdown")
	case msg := <-errChan:
		log.Error("error: ", msg)
	}

	db.Close()
	os.Exit(1)
}
