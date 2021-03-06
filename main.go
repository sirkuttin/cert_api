package main

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/dgraph-io/badger"
	"github.com/sirkuttin/acme-cert-api/api"
	"os"
	"os/signal"
	"syscall"
)

var log = logrus.New()

func init() {
	log.SetLevel(logrus.DebugLevel)
}

func main() {

	db, err := createDB()
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

func createDB() (db *badger.DB, err error) {
	opts := badger.DefaultOptions
	opts.Dir = "app_db"
	opts.ValueDir = "app_db"
	return badger.Open(opts)
}
