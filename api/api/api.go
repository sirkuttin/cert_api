package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	log logrus.Logger
)

func Start(logger *logrus.Logger) {

	log = *logger

	log.Info("Starting API")

	router := mux.NewRouter()

	router.HandleFunc("/cert/{domain}", getCert()).Methods("GET")

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
