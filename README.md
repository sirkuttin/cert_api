# Cert API

This is a simple api with one endpoint that accepts GET and POST requests.

GET /cert/{domain}
Returns a json response that includes the cert string. If the cert does not exist, this endpoint will create it with default expiration.
Upon a certificate hitting it's expiration, the api will autogenerate and return a new cert string.

POST /cert/{domain}
Pass in a json request with the certificates time to live in minutes. ttl_minutes defines how long until the cert expires and a new cert is auto-generated.
Returns a json response that includes the cert string.

## Prerequisites

Install Golang


## API Postman Collection
https://documenter.getpostman.com/view/52582/RWTprGXy
