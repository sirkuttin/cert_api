package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"github.com/dgraph-io/badger"
	"io/ioutil"
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/pkg/errors"
)

type Certificate struct {
	Domain            string        `json:"-"`
	Cert              string        `json:"cert"`
	TimeToLiveMinutes time.Duration `json:"ttl_minutes"`
	Expiration        time.Time     `json:"expiration"`
}

func addCertificate() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		addDelay := time.NewTimer(10 * time.Second)

		cert, err := createCertificateFromRequest(request)
		if err != nil {
			http.Error(responseWriter, err.Error(), 400)
		}

		cert, err = upsertCertificate(cert)
		if err != nil {
			http.Error(responseWriter, err.Error(), 500)
			return
		}

		certBytes, _ := json.Marshal(cert)
		fmt.Fprint(responseWriter, string(certBytes))

		<- addDelay.C
	}
}

func createCertificateFromRequest(request *http.Request) (cert Certificate, err error) {
	bodyBytes, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyBytes, &cert)
	if err != nil {
		return
	}
	domain := mux.Vars(request)["domain"]
	if domain == "" {
		err = errors.New("domain cannot be empty")
		return
	}
	cert.Domain = domain
	return
}

func getCertificate() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {

		addDelay := time.NewTimer(10 * time.Second)

		domain := mux.Vars(request)["domain"]
		if domain == "" {
			fmt.Fprintf(responseWriter, "domain cannot be empty")
			return
		}
		log.Debug("domain = ", domain)

		db.View(func(transaction *badger.Txn) error {
			item, err := transaction.Get([]byte(domain))
			if err != nil && item != nil {
				return err
			}
			if item == nil {
				cert, _ := upsertCertificate(Certificate{Domain: domain})
				certBytes, _ := json.Marshal(cert)
				fmt.Fprint(responseWriter, string(certBytes))
				return nil
			}

			val, err := item.Value()
			if err != nil {
				return err
			}

			var cert Certificate;
			err = json.Unmarshal(val, &cert)
			if err != nil {
				http.Error(responseWriter, err.Error(), 500)
				return err
			}

			cert.Domain = domain

			if cert.Expiration.Before(time.Now()) {
				cert, err = upsertCertificate(cert)
				if err != nil {
					http.Error(responseWriter, err.Error(), 500)
					return nil
				}
			}
			certBytes, _ := json.Marshal(cert)
			fmt.Fprint(responseWriter, string(certBytes))
			return nil
		})

		<- addDelay.C
	}
}

func upsertCertificate(cert Certificate) (Certificate, error) {

	if cert.TimeToLiveMinutes.Minutes() == 0 {
		cert.TimeToLiveMinutes = 5
	}
	cert.Expiration = time.Now().Add(cert.TimeToLiveMinutes * time.Minute)
	cert.Cert = uuid.NewV1().String()
	certBytes, _ := json.Marshal(cert)
	err := db.Update(func(transaction *badger.Txn) error {
		err := transaction.Set([]byte(cert.Domain), certBytes)
		return err
	})
	return cert, err
}
