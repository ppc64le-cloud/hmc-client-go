package hmc

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type BasicAuthenticator struct {
	UserName string
	Password string
	BaseURL  string
	Client   *http.Client
}

type Attributes struct {
	XMLns         string `xml:"xmlns,attr"`
	SchemaVersion string `xml:"schemaVersion,attr"`
}

var attrib = Attributes{"http://www.ibm.com/xmlns/systems/power/firmware/web/mc/2012_10/", "V1_0"}

func (a *BasicAuthenticator) Validate() error {
	if a.UserName == "" || a.Password == "" {
		return fmt.Errorf("username or password can't be empty")
	}

	if a.BaseURL == "" {
		return fmt.Errorf("username or password can't be empty")
	}
	return nil
}

func (a *BasicAuthenticator) GetBaseURL() string {
	return a.BaseURL
}

func (a *BasicAuthenticator) SetClient(client *http.Client) {
	a.Client = client
}

func (a *BasicAuthenticator) Authenticate(request *http.Request) error {
	if err := a.Validate(); err != nil {
		return err
	}

	u, err := url.Parse(a.BaseURL)
	if err != nil {
		return err
	}

	// Will return if already cookies present for the URL, no need to call login function.
	if cookies := a.Client.Jar.Cookies(u); len(cookies) > 0 {
		return nil
	}

	if err := a.login(); err != nil {
		return fmt.Errorf("login failed with error: %v", err)
	}

	return nil
}

func (a *BasicAuthenticator) login() error {
	l := struct {
		Attributes
		XMLName  xml.Name `xml:"LogonRequest"`
		UserID   string   `xml:"UserID"`
		Password string   `xml:"Password"`
	}{
		Attributes: attrib,
		UserID:     a.UserName,
		Password:   a.Password,
	}

	body, err := xml.Marshal(l)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, a.GetBaseURL()+"/rest/api/web/Logon", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/vnd.ibm.powervm.web+xml; type=LogonRequest")

	res, err := a.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err := json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status: %s", res.Status)
	}
	return nil
}
