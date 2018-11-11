package gotoon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	authorizeURL = "https://api.toon.eu/authorize"
	tokenURL     = "https://api.toon.eu/token"
	apiBaseURL   = "https://api.toon.eu/toon/v3"
)

// jsonTime defines customized JSON marshal and unmarshal functions
// for converting timestamp into Time struct.
type jsonTime time.Time

// MarshalJSON marshals Time struct into timestamp integer in milliseconds.
func (t jsonTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

// UnmarshalJSON unmarshals timestamp integer in milliseconds into Time struct.
func (t *jsonTime) UnmarshalJSON(s []byte) (err error) {
	q, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q/1000, 0)
	return
}

func (t jsonTime) String() string { return time.Time(t).String() }

// jsonBool defines customized JSON unmarshal function for converting
// integer into boolean.
type jsonBool bool

// UnmarshalJSON converts input string or integer into a boolean value.
func (b *jsonBool) UnmarshalJSON(s []byte) (err error) {
	bs := string(s)
	if bs == "0" || bs == "false" {
		*b = false
	} else if bs == "1" || bs == "true" {
		*b = true
	} else {
		err = fmt.Errorf("Cannot unmarshal value to boolean: %s", bs)
	}
	return
}

// token holds the data structure of the Toon API access token.
type token struct {
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

// ThermostatStates holds the data structure of the last states retrieved from
// the getStatus interface of the Toon API.
type ThermostatStates struct {
	State                  []ThermostatState `json:"state"`
	LastUpdatedFromDisplay jsonTime          `json:"lastUpdatedFromDisplay,int"`
}

// ThermostatState holds the data structure of a state retrieved from the getStatus
// interface of the Toon API.
type ThermostatState struct {
	ID        int `json:"id"`
	TempValue int `json:"tempValue"`
	Dhw       int `json:"dhw"`
}

// ThermostatInfo holds the data structure of the thermostat information retrieved
// from the getStatus interface of the Toon API.
type ThermostatInfo struct {
	CurrentSetPoint        int      `json:"currentSetpoint"`
	CurrentDisplayTemp     int      `json:"currentDisplayTemp"`
	ProgramState           int      `json:"programState"`
	ActiveState            int      `json:"activeState"`
	NextProgram            int      `json:"nextProgram"`
	NextState              int      `json:"nextState"`
	NextTime               int      `json:"nextTime"`
	NextSetPoint           int      `json:"nextSetpoint"`
	ErrorFound             int      `json:"errorFound"`
	BoilerModuleConnected  int      `json:"boilerModuleConnected"`
	RealSetPoint           int      `json:"realSetpoint"`
	BurnerInfo             string   `json:"burnerInfo"`
	OtCommError            string   `json:"otCommError"`
	CurrentModulationLevel int      `json:"currentModulationLevel"`
	HaveOTBoiler           int      `json:"haveOTBoiler"`
	LastUpdatedFromDisplay jsonTime `json:"lastUpdatedFromDisplay,int"`
}

// PowerUsage holds the data structure of the current power consumption retrieved from
// the getStatus interface of the Toon API.
type PowerUsage struct {
	Value                  int      `json:"value"`
	DayCost                float32  `json:"dayCost,int"`
	ValueProduced          int      `json:"valueProduced"`
	DayCostProduced        int      `json:"dayCostProduced"`
	ValueSolar             int      `json:"valueSolar"`
	MaxSolar               int      `json:"maxSolar"`
	DayCostSolar           int      `json:"dayCostSolar"`
	AvgSolarValue          int      `json:"avgSolarValue"`
	AvgValue               float32  `json:"avgValue"`
	AvgDayValue            float32  `json:"avgDayValue"`
	AvgProduValue          int      `json:"avgProduValue"`
	AvgDayProduValue       int      `json:"avgDayProduValue"`
	MeterReading           int      `json:"meterReading"`
	MeterReadingLow        int      `json:"meterReadingLow"`
	MeterReadingProdu      int      `json:"meterReadingProdu"`
	MeterReadingLowProdu   int      `json:"meterReadingLowProdu"`
	DayUsage               int      `json:"dayUsage"`
	DayLowUsage            int      `json:"dayLowUsage"`
	TodayLowestUsage       int      `json:"todayLowestUsage"`
	IsSmart                jsonBool `json:"isSmart,int"`
	LowestDayValue         int      `json:"lowestDayValue"`
	SolarProducedToday     int      `json:"solarProducedToday"`
	LastUpdatedFromDisplay jsonTime `json:"lastUpdatedFromDisplay,int"`
}

// GasUsage holds the data structure of the current gas consumption retrieved from
// the getStatus interface of the Toon API.
type GasUsage struct {
	Value                  int      `json:"value"`
	DayCost                float32  `json:"dayCost"`
	AvgValue               float32  `json:"avgValue"`
	MeterReading           int      `json:"meterReading"`
	AvgDayValue            float32  `json:"avgDayValue"`
	DayUsage               int      `json:"dayUsage"`
	IsSmart                jsonBool `json:"isSmart,int"`
	LastUpdatedFromDisplay jsonTime `json:"lastUpdatedFromDisplay,int"`
}

// Status holds the main data structure of the current Toon device status retrieved
// from the getStatus interface of the Toon API.
type Status struct {
	ThermostatStates      ThermostatStates `json:"thermostatStates"`
	ThermostatInfo        ThermostatInfo   `json:"thermostatInfo"`
	PowerUsage            PowerUsage       `json:"powerUsage"`
	GasUsage              GasUsage         `json:"gasUsage"`
	LastUpdateFromDisplay jsonTime         `json:"lastUpdateFromDisplay,int"`
}

// FlowDataPoint holds the data structure of the consumption data points.
type FlowDataPoint struct {
	Timestamp jsonTime `json:"timestamp,int"`
	Unit      string   `json:"unit"`
	Value     float32  `json:"value"`
}

// FlowData holds the data structure of the consumption data.
type FlowData struct {
	Hours  []FlowDataPoint `json:"hours"`
	Days   []FlowDataPoint `json:"days"`
	Weeks  []FlowDataPoint `json:"weeks"`
	Months []FlowDataPoint `json:"months"`
	Years  []FlowDataPoint `json:"years"`
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
	accessToken token
}

// getAccessToken authorise the user to get the access token for retriving data
// from the Toon device.
func (t *Toon) getAccessToken() (err error) {

	c := newHTTPSClient()

	// step 1: call https://api.toon.eu/authorize (optionally?)
	//         with input: client_id, response_type=code, redirect_url=http://127.0.0.1, tenant_id
	//         This step doesn't seem to be necessary.  Comment it out for the moment.
	// v := url.Values{}
	// v.Set("client_id", t.ConsumerKey)
	// v.Set("response_type", "code")
	// v.Set("redirect_url", "http://127.0.0.1")
	// v.Set("tenant_id", t.TenantID)

	// _, err = c.Get(authorizeURL + v.Encode())
	// if err != nil {
	// 	return
	// }

	// step 2: call https://api.toon.eu/authorize/legacy to get "code" from the returned HTTP header
	//         with input: client_id, tenant_id, username, password, response_type=code,
	//         state='', scope=''
	v := url.Values{}
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

	// unmarshal response body to Token struct
	if err = json.Unmarshal(bodyBytes, &(t.accessToken)); err != nil {
		return
	}

	// derive ExpiresAt = tnow + (ExpiresIn - 180)s
	t.accessToken.ExpiresAt = tnow.Add(time.Second * time.Duration(t.accessToken.ExpiresIn-180))
	t.accessToken.RefreshTokenExpiresAt = tnow.Add(time.Second * time.Duration(t.accessToken.RefreshTokenExpiresIn-180))

	return
}

func (t *Toon) refreshAccessToken() (err error) {
	// TODO: try to refresh the token using the token refresh function
	return
}

func (t *Toon) hasValidToken() (isValid bool) {

	// the accessToken is not set
	if &(t.accessToken) == nil {
		isValid = false
		return
	}

	// the refresh token has been expired
	if t.accessToken.RefreshTokenExpiresAt.Before(time.Now()) {
		isValid = false
		return
	}

	// the token has expired; but we can try to renew the token
	if t.accessToken.ExpiresAt.Before(time.Now()) {
		// given the refresh token is still valid, try refreshing the access token.
		if err := t.refreshAccessToken(); err != nil {
			isValid = false
			return
		}
	}

	// finally check whether the current/refreshed access token is valid.
	isValid = t.accessToken.ExpiresAt.Before(time.Now())
	return
}

// GetAgreements gets identifier information of accessible Toon devices.
func (t *Toon) GetAgreements() (agreements []Agreement, err error) {

	if !t.hasValidToken() {
		if err = t.getAccessToken(); err != nil {
			return
		}
	}

	c := newHTTPSClient()
	req, err := http.NewRequest("GET", apiBaseURL+"/agreements", nil)
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

	if res.StatusCode != 200 {
		err = fmt.Errorf("Cannot get agreements: %s", string(bodyBytes))
		return
	}

	err = json.Unmarshal(bodyBytes, &agreements)
	return
}

// GetStatus returns current information about the thermostat status,
// and usages of electricity and gas of a Toon device identified by the
// given Agreement.
//
// The information is retrieved via the Toon API endpoint:
// https://api.toon.eu/toon/v3/{agreement.AgreementID}/status
func (t *Toon) GetStatus(agreement Agreement) (status Status, err error) {

	if &(agreement.AgreementID) == nil {
		err = fmt.Errorf("Invalid agreement: %+v", agreement)
		return
	}

	if !t.hasValidToken() {
		if err = t.getAccessToken(); err != nil {
			return
		}
	}

	c := newHTTPSClient()
	for {
		var req *http.Request
		req, err = http.NewRequest("GET", apiBaseURL+"/"+agreement.AgreementID+"/status", nil)
		if err != nil {
			break
		}
		req.Header.Set("authorization", "Bearer "+t.accessToken.AccessToken)
		req.Header.Set("accept", "application/json")
		req.Header.Set("cache-control", "no-cache")
		req.Header.Set("content-type", "application/json")

		var res *http.Response
		res, err = c.Do(req)
		if err != nil {
			break
		}

		// read the response body
		var bodyBytes []byte
		bodyBytes, err = ioutil.ReadAll(res.Body)
		if err != nil {
			break
		}

		// the data is retrieved. Unmarshal the JSON document into
		// the Status data structure.
		if res.StatusCode == 200 {
			err = json.Unmarshal(bodyBytes, &status)
			break
		}

		// the server may accept the request, and require the client to
		// retrieve the information later.  Looks like the information
		// of device is retrieved on demand??
		// In this case, the client receive status 202, and we need to
		// send the same request again with the same access token until
		// we got the data.
		if res.StatusCode != 202 {
			err = fmt.Errorf("Error getting status: %s", string(bodyBytes))
			break
		}
	}
	return
}

// GetGasFlow retrieves gas consumption information from a given Toon device
// for a given time period in 5-minute intervals.  The Toon device is referred
// by the agreement parameter, and the time period is indicated by the from-
// and toTime parameters.
//
// The information is retrieved via the Toon API endpoint:
// https://api.toon.eu/toon/v3/{agreement.AgreementID}/consumption/gas/flows
func (t *Toon) GetGasFlow(agreement Agreement, fromTime, toTime time.Time) (flow FlowData, err error) {
	if &(agreement.AgreementID) == nil {
		err = fmt.Errorf("Invalid agreement: %+v", agreement)
		return
	}

	if !t.hasValidToken() {
		if err = t.getAccessToken(); err != nil {
			return
		}
	}

	c := newHTTPSClient()
	var req *http.Request
	var res *http.Response
	var bodyBytes []byte

	// compose API endpoint
	endpointURL := apiBaseURL + "/" + agreement.AgreementID + "/consumption/gas/flows"
	req, err = http.NewRequest("GET", endpointURL, nil)
	if err != nil {
		return
	}

	// add query parameter values to the request
	v := req.URL.Query()
	if (time.Time{}) != fromTime {
		v.Add("fromTime", fmt.Sprintf("%d", 1000*fromTime.Unix()))
	}
	if (time.Time{}) != toTime {
		v.Add("toTime", fmt.Sprintf("%d", 1000*toTime.Unix()))
	}
	req.URL.RawQuery = v.Encode()

	req.Header.Set("authorization", "Bearer "+t.accessToken.AccessToken)
	req.Header.Set("accept", "application/json")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/json")

	res, err = c.Do(req)
	if err != nil {
		return
	}

	// read the response body
	bodyBytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		err = fmt.Errorf("Error getting gas consumption flow: %s", string(bodyBytes))
	}

	err = json.Unmarshal(bodyBytes, &flow)

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
