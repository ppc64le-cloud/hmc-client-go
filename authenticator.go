package hmc

import "net/http"

type Authenticator interface {
	Authenticate(request *http.Request) error
	Validate() error
	GetBaseURL() string
	SetClient(client *http.Client)
}
