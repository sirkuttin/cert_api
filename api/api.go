package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/dgraph-io/badger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	log logrus.Logger
	db  *badger.DB
)

func Start(logger *logrus.Logger, database *badger.DB) {

	log = *logger
	log.Info("Starting API")
	db = database

	router := mux.NewRouter()
	router.HandleFunc("/cert/{domain}", getCertificate()).Methods("GET")
	router.HandleFunc("/cert/{domain}", addCertificate()).Methods("POST")

	err := http.ListenAndServe(":8000", handlers.CORS(createCorsOptions()...)(router))
	if err != nil {
		panic(err.Error())
	}
}

func createCorsOptions() (corsOptions []handlers.CORSOption) {
	corsOptions = append(corsOptions, handlers.AllowedHeaders([]string{"X-Requested-With"}))
	corsOptions = append(corsOptions, handlers.AllowedOrigins([]string{"*"}))
	corsOptions = append(corsOptions, handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}))
	return
}
