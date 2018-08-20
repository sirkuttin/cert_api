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

type Cert struct {
	Domain            string        `json:"-"`
	Cert              string        `json:"cert"`
	TimeToLiveMinutes time.Duration `json:"ttl_minutes"`
	Expiration        time.Time     `json:"expiration"`
}

func addCert() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {

		cert, err := createCertFromRequest(request)
		if err != nil {
			http.Error(responseWriter, err.Error(), 400)
		}

		cert, err = upsertCert(cert)
		if err != nil {
			http.Error(responseWriter, err.Error(), 500)
			return
		}

		certBytes, _ := json.Marshal(cert)
		fmt.Fprint(responseWriter, string(certBytes))
	}
}

func createCertFromRequest(request *http.Request) (cert Cert, err error) {
	bodyBytes, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		//http.Error(responseWriter, err.Error(), 500)
		return
	}
	// create cert struct with ttl
	err = json.Unmarshal(bodyBytes, &cert)
	if err != nil {
		//http.Error(responseWriter, err.Error(), 500)
		return
	}

	domain := mux.Vars(request)["domain"]
	if domain == "" {
		err = errors.New("domain cannot be empty")
		return //http.Error(responseWriter, "domain cannot be empty", 400)
	}
	cert.Domain = domain
	return
}

func getCert() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		domain := mux.Vars(request)["domain"]
		if domain == "" {
			fmt.Fprintf(responseWriter, "domain cannot be empty")
			return
		}

		log.Debug("domain = ", domain)
		//time.Sleep(10 * time.Second) /////////////////////////////////////////////////////////////////////

		db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(domain))
			if err != nil && item != nil {
				return err
			}

			if item == nil {
				cert, _ := upsertCert(Cert{Domain:domain})
				certBytes, _ := json.Marshal(cert)
				fmt.Fprint(responseWriter, string(certBytes))
				return nil
			}

			val, err := item.Value()
			if err != nil {
				return err
			}

			var cert Cert
			err = json.Unmarshal(val, &cert)
			if err != nil {
				http.Error(responseWriter, err.Error(), 500)
				return err
			}

			cert.Domain = domain

			if (cert.Expiration.Before(time.Now())) {
				cert, err = upsertCert(cert)
				if err != nil {
					http.Error(responseWriter, err.Error(), 500)
					return nil
				}
			}
			certBytes, _ := json.Marshal(cert)
			fmt.Fprint(responseWriter, string(certBytes))
			return nil
		})
	}
}

func upsertCert(cert Cert) (Cert, error) {

	if cert.TimeToLiveMinutes.Minutes() == 0 {
		cert.TimeToLiveMinutes = 5
	}
	cert.Expiration = time.Now().Add(cert.TimeToLiveMinutes * time.Minute)

	//add uuid cert string to the cert struct
	certString, _ := uuid.NewV1()
	cert.Cert = certString.String()

	if cert.Domain == "" {
		return cert, errors.New("domain cannot be empty")
	}

	log.Debug("domain = ", cert.Domain)
	//time.Sleep(10 * time.Second) //////////////////////////////////////////////////////

	certBytes, _ := json.Marshal(cert)
	// create key = domain value = json of newly created uuid and expiration date
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(cert.Domain), certBytes)
		return err
	})
	return cert, err
}
