package gotoon_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hurngchunlee/gotoon"
)

var toon gotoon.Toon

func init() {
	toon = gotoon.Toon{
		Username:       os.Getenv("TOONAPI_TEST_USERNAME"),
		Password:       os.Getenv("TOONAPI_TEST_PASSWORD"),
		TenantID:       "eneco",
		ConsumerKey:    os.Getenv("TOONAPI_TEST_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TOONAPI_TEST_CONSUMER_SECRET"),
	}
}

func TestGetAgreements(t *testing.T) {

	agreements, err := toon.GetAgreements()
	if err != nil {
		t.Errorf("Fail getting agreements: %+v\n", err)
	}

	for _, agreement := range agreements {
		t.Logf("%+v", agreement)
	}
}

func TestGetStatus(t *testing.T) {

	agreements, err := toon.GetAgreements()
	if err != nil {
		t.Errorf("Fail getting agreements: %+v\n", err)
	}

	for _, agreement := range agreements {
		status, err := toon.GetStatus(agreement)
		if err != nil {
			t.Errorf("%s: fail getting status - %+v\n", agreement.AgreementID, err)
		}
		t.Logf("%s: %+v", agreement.AgreementID, status)
	}
}

func TestGetGasFlow(t *testing.T) {
	agreements, err := toon.GetAgreements()
	if err != nil {
		t.Errorf("Fail getting agreements: %+v\n", err)
	}

	for _, agreement := range agreements {
		// retrieve gas consumption of the last 30 minutes in 5-minute intervals
		var fromTime time.Time
		var toTime time.Time
		toTime = time.Now()
		fromTime = (toTime).Add(time.Duration(-30) * time.Minute)

		flow, err := toon.GetGasFlow(agreement, fromTime, toTime)
		if err != nil {
			t.Errorf("%s: fail getting flow - %+v\n", agreement.AgreementID, err)
		}
		t.Logf("Gas flow %s: %+v", agreement.AgreementID, flow)
	}
}

// The code below shows how to get Toon device agreements.
func ExampleToon_GetAgreements() {
	toon := gotoon.Toon{
		Username:       "myEnecoUsername",
		Password:       "myEnecoPassword",
		TenantID:       "eneco",
		ConsumerKey:    "ToonAPIConsumerKey",
		ConsumerSecret: "ToonAPIConsumerSecret",
	}

	agreements, err := toon.GetAgreements()
	if err != nil {
		fmt.Printf("Fail getting agreements: %+v\n", err)
	}

	for _, agreement := range agreements {
		fmt.Printf("%+v", agreement)
	}
}

// The code below shows how to get current status of Toon devices.
func ExampleToon_GetStatus() {
	toon := gotoon.Toon{
		Username:       "myEnecoUsername",
		Password:       "myEnecoPassword",
		TenantID:       "eneco",
		ConsumerKey:    "ToonAPIConsumerKey",
		ConsumerSecret: "ToonAPIConsumerSecret",
	}

	agreements, err := toon.GetAgreements()
	if err != nil {
		fmt.Printf("Fail getting agreements: %+v\n", err)
	}

	for _, agreement := range agreements {
		status, _ := toon.GetStatus(agreement)
		fmt.Printf("%+v", status)
	}
}

// The code below shows how to get gas consumption during the last 30 minutes,
// in 5-minute intervals.
func ExampleToon_GetGasFlow() {
	toon := gotoon.Toon{
		Username:       "myEnecoUsername",
		Password:       "myEnecoPassword",
		TenantID:       "eneco",
		ConsumerKey:    "ToonAPIConsumerKey",
		ConsumerSecret: "ToonAPIConsumerSecret",
	}

	agreements, err := toon.GetAgreements()
	if err != nil {
		fmt.Printf("Fail getting agreements: %+v\n", err)
	}

	for _, agreement := range agreements {
		var fromTime time.Time
		var toTime time.Time
		// Set period of the last 30 minutes by defining the fromTime and toTime parameters.
		// You may omit it for the default period of the last 1 hour.
		toTime = time.Now()
		fromTime = toTime.Add(time.Duration(-30) * time.Minute)

		flow, _ := toon.GetGasFlow(agreement, fromTime, toTime)
		fmt.Printf("Gas flow %s: %+v\n", agreement.AgreementID, flow)
	}
}
