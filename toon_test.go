package toonapilib

import (
	"os"
	"testing"
)

var toon Toon

func init() {
	toon = Toon{
		Username:       os.Getenv("TOONAPI_TEST_USERNAME"),
		Password:       os.Getenv("TOONAPI_TEST_PASSWORD"),
		TenantID:       "eneco",
		ConsumerKey:    os.Getenv("TOONAPI_TEST_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TOONAPI_TEST_CONSUMER_SECRET"),
	}
}

func TestGetAccessToken(t *testing.T) {
	err := toon.getAccessToken()
	if err != nil {
		t.Errorf("Fail getting access token: %+v\n", err)
	}

	t.Logf("%+v", toon.accessToken)
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
