package hmc

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type RoundTripper struct {
	Proxied http.RoundTripper
}

func (rt RoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	// Log the request body
	log.WithFields(log.Fields{
		"URL":     req.URL.String(),
		"headers": req.Header,
	}).Debugln("request")
	// Send the request, get the response (or the error)
	res, e = rt.Proxied.RoundTrip(req)

	// Handle the result.
	if e != nil {
		log.Debugf("Error: %v", e)
	} else {
		log.Debugf("Received %v response", res.Status)
	}

	return
}
