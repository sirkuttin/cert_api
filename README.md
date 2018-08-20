# Cert API

This is a simple api with one endpoint that accepts GET and POST requests.

GET /cert/{domain}
Returns a json response that includes the cert string. If the cert does not exist, this endpoint will create it with default expiration.
Upon a certificate hitting it's expiration, the api will autogenerate and return a new cert string.

POST /cert/{domain}
Pass in a json request with the certificates time to live in minutes. ttl_minutes defines how long until the cert expires and a new cert is auto-generated.
Returns a json response that includes the cert string.

## How To Run
Pick up the latest release from the release tab.
#### Linux
```
chmod a+x acme-cert-api
./acme-cert-api
```
#### Windows
open cmd, type in the file name, press enter

## API Postman Collection
https://documenter.getpostman.com/view/52582/RWTprGXy

## How To Run From Source
* Install Golang
* Install Dep
```
mkdir -p $GOPATH/src/github.com/sirkuttin
cd $GOPATH/src/github.com/sirkuttin/
git clone https://github.com/sirkuttin/cert_api.git
cd cert_api
dep ensure
go run main.go
```


