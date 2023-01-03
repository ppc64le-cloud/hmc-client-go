package hmc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Client contains parameters for configuring the SDK.
type Client struct {
	// baseURL used for the hmc REST API Calls
	baseURL string

	body []byte

	// Authenticator for the client
	Authenticator Authenticator

	// apiKey used to talk ManageIQ
	apiKey string

	// No need to set -- for testing only
	HTTPClient *http.Client
}

func (c *Client) GetBaseURL() string {
	return c.baseURL
}
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

type ClientParams struct {
	LogLevel log.Level
	Insecure bool
}

func NewClient(authenticator Authenticator, param ClientParams) *Client {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(param.LogLevel)
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: param.Insecure}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	c := &http.Client{
		Timeout:   time.Minute,
		Transport: RoundTripper{tr},
		Jar:       jar,
	}
	authenticator.SetClient(c)

	return &Client{
		baseURL:       authenticator.GetBaseURL() + "/rest/api/uom",
		Authenticator: authenticator,
		HTTPClient:    c,
	}
}

// DetailedResponse holds the response information received from the server.
type DetailedResponse struct {

	// The HTTP status code associated with the response.
	StatusCode int

	// The HTTP headers contained in the response.
	Headers http.Header

	// Result - this field will contain the result of the operation (obtained from the response body).
	Result *Feed

	// RawResult field will contain the raw response body.
	RawResult []byte
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{}
}

func (c *Client) sendRequest(req *http.Request) (*Feed, *DetailedResponse, error) {
	if err := c.Authenticator.Authenticate(req); err != nil {
		return nil, nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return nil, nil, errors.New(errRes.Message)
		}

		return nil, nil, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	detailedResponse := &DetailedResponse{
		StatusCode: res.StatusCode,
		Headers:    res.Header,
		Result:     &Feed{},
		RawResult:  []byte(trimEmptyLines(body)),
	}

	if err := xml.Unmarshal(detailedResponse.RawResult, detailedResponse.Result); err != nil {
		return nil, nil, err
	}

	return detailedResponse.Result, detailedResponse, nil
}

func trimEmptyLines(b []byte) string {
	strs := strings.Split(string(b), "\n")
	str := ""
	for _, s := range strs {
		if len(strings.TrimSpace(s)) == 0 {
			continue
		}
		str += s + "\n"
	}
	str = strings.TrimSuffix(str, "\n")

	return str
}

func (c *Client) GET() (*Feed, *DetailedResponse, error) {
	req, err := http.NewRequest(http.MethodGet, c.GetBaseURL(), nil)
	req.Header.Set("Accept", "application/atom+xml; charset=UTF-8")
	if err != nil {
		return nil, nil, err
	}
	return c.sendRequest(req)
}

func (c *Client) POST() (*Feed, *DetailedResponse, error) {
	req, err := http.NewRequest(http.MethodPost, c.GetBaseURL(), nil)
	if err != nil {
		return nil, nil, err
	}
	return c.sendRequest(req)
}

func (c *Client) PUT() (*Feed, *DetailedResponse, error) {
	req, err := http.NewRequest(http.MethodPut, c.GetBaseURL(), bytes.NewReader(c.body))
	if err != nil {
		return nil, nil, err
	}
	return c.sendRequest(req)
}

func (c *Client) DELETE() (*Feed, *DetailedResponse, error) {
	req, err := http.NewRequest(http.MethodDelete, c.GetBaseURL(), nil)
	if err != nil {
		return nil, nil, err
	}
	return c.sendRequest(req)
}

type Resource interface {
	GET() interface{}
	DELETE() error
	POST() interface{}
}
