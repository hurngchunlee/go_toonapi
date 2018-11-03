package gotoon_test

import (
	"fmt"
	"gotoon"
	"os"
	"testing"
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

// This function shows how of getting Toon API agreements.
func ExampleGetAgreements() {
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
