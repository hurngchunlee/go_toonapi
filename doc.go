/*

Package gotoon is a Go library for accessing the Toon (https://www.toon.eu) smart thermostat,
using the Toon RESTful API (https://developer.toon.eu).

Installation

    go get -u github.com/hurngchunlee/gotoon

Usage

A small example:

    package main

    import (
        "fmt"
        "github.com/hurngchunlee/gotoon"
    )

    func main() {

        // Initialize Toon with authentication credentials.
	toon := gotoon.Toon{
		Username:       "myEnecoUsername",
		Password:       "myEnecoPassword",
		TenantID:       "eneco",
		ConsumerKey:    "ToonAPIConsumerKey",
		ConsumerSecret: "ToonAPIConsumerSecret",
	}

        // Call method to retrieve information or interact with the Toon device.
        // Example 1: Get agreementIds.
	agreements, err := toon.GetAgreements()

    }

*/
package gotoon
