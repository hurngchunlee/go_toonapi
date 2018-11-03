package gotoon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	authorizeURL = "https://api.toon.eu/authorize"
	tokenURL     = "https://api.toon.eu/token"
	agreementURL = "https://api.toon.eu/toon/v3/agreements"
)

// Token holds the data structure of the Toon API access token.
type Token struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in,string"`
	ExpiresAt             time.Time
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in,string"`
	RefreshTokenExpiresAt time.Time
}

// Agreement holds the data structure of the Toon API agreement. See https://developer.toon.eu/api-intro.
type Agreement struct {
	AgreementID            string `json:"agreementId"`
	AgreementIDChecksum    string `json:"agreementIdChecksum"`
	HeatingType            string `json:"heatingType"`
	DisplayCommonName      string `json:"displayCommonName"`
	DisplayHardwareVersion string `json:"displayHardwareVersion"`
	DisplaySoftwareVersion string `json:"displaySoftwareVersion"`
	IsToonSolar            bool   `json:"isToonSolar"`
	IsToonly               bool   `json:"isToonly"`
}

// Toon provides interface to access and retrieve data from the Toon device,
// using the Toon RESTful APIs, see https://developer.toon.eu.
type Toon struct {
	// Username is the Tenant account name for the Tenant (e.g. the Mijn Eneco account)
	Username string
	// Password is the password for the Tenant account (e.g. the Mijn Eneco password)
	Password string
	// TenantID is the tenant ID (e.g. eneco, viesgo)
	TenantID string
	// ConsumerKey is the consumer key of the Toon API, see https://developer.toon.eu/authentication
	ConsumerKey string
	// ConsumerSecret is the consumer secret of the Toon API, see https://developer.toon.eu/authentication
	ConsumerSecret string
	// accessToken is the current Toon API access token, see https://developer.toon.eu/authentication
	accessToken Token
}

// getAccessToken authorise the user to get the access token for retriving data
// from the Toon device.
func (t *Toon) getAccessToken() (err error) {

	c := newHTTPSClient()

	// step 1: call https://api.toon.eu/authorize (optionally?)
	//         with input: client_id, response_type=code, redirect_url=http://127.0.0.1, tenant_id
	v := url.Values{}
	v.Set("client_id", t.ConsumerKey)
	v.Set("response_type", "code")
	v.Set("redirect_url", "http://127.0.0.1")
	v.Set("tenant_id", t.TenantID)

	_, err = c.Get(authorizeURL + v.Encode())
	if err != nil {
		return
	}

	// step 2: call https://api.toon.eu/authorize/legacy to get "code" from the returned HTTP header
	//         with input: client_id, tenant_id, username, password, response_type=code,
	//         state='', scope=''
	v = url.Values{}
	v.Set("client_id", t.ConsumerKey)
	v.Set("tenant_id", t.TenantID)
	v.Set("username", t.Username)
	v.Set("password", t.Password)
	v.Set("response_type", "code")
	v.Set("state", "")
	v.Set("scope", "")

	// disable http redirect
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	r, err := c.PostForm(authorizeURL+"/legacy", v)
	if err != nil {
		return
	}
	if r.StatusCode != 302 {
		err = errors.New("invalid consumer key")
		return
	}

	// parsing Location header attribute to get value of the "code"
	u, err := url.Parse(r.Header.Get("Location"))
	if err != nil {
		return
	}
	code := u.Query().Get("code")
	if code == "" {
		err = fmt.Errorf("fail extracting code, header: +%v", r.Header)
	}
	log.Printf("code: " + code)

	// step 3: call https://api.toon.eu/token to get the access token
	v = url.Values{}
	v.Set("client_id", t.ConsumerKey)
	v.Set("client_secret", t.ConsumerSecret)
	v.Set("grant_type", "authorization_code")
	v.Set("code", code)

	// current time
	tnow := time.Now()
	r, err = c.PostForm(tokenURL, v)
	if err != nil {
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	log.Printf(string(bodyBytes))

	// unmarshal response body to Token struct
	if err = json.Unmarshal(bodyBytes, &(t.accessToken)); err != nil {
		return
	}

	// derive ExpiresAt = tnow + (ExpiresIn - 180)s
	t.accessToken.ExpiresAt = tnow.Add(time.Second * time.Duration(t.accessToken.ExpiresIn-180))
	t.accessToken.RefreshTokenExpiresAt = tnow.Add(time.Second * time.Duration(t.accessToken.RefreshTokenExpiresIn-180))

	return
}

func (t *Toon) hasValidToken() (isValid bool) {

	// the accessToken is not set
	if &(t.accessToken) == nil {
		isValid = false
		return
	}

	// the refresh token has been expired
	if t.accessToken.RefreshTokenExpiresAt.After(time.Now()) {
		isValid = false
		return
	}

	// the token has expired; but we can try to renew the token
	if t.accessToken.ExpiresAt.After(time.Now()) {
		// TODO: try to refresh the token using the token refresh function
		isValid = false
		return
	}

	isValid = true
	return
}

// GetAgreements gets identifier information of accessible Toon devices.
func (t *Toon) GetAgreements() (agreements []Agreement, err error) {

	if t.hasValidToken() {
		if err = t.getAccessToken(); err != nil {
			return
		}
	}

	c := newHTTPSClient()
	req, err := http.NewRequest("GET", agreementURL, nil)
	if err != nil {
		return
	}
	req.Header.Set("authorization", "Bearer "+t.accessToken.AccessToken)
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bodyBytes, &agreements)
	return
}

// internal utility functions
func newHTTPSClient() (client *http.Client) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	client = &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	return
}
