package gotoon_test

import (
	"fmt"
	"os"
	"testing"

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

// The code below shows how of getting Toon API agreements.
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
