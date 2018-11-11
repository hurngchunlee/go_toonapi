/*

Package gotoon is a Go library for accessing the Toon (https://www.toon.eu) smart thermostat,
using the Toon API (https://developer.toon.eu).

Installation

Use the following command to install the library.

    go get -u github.com/hurngchunlee/gotoon

Usage

A small example of retriving the current status and device information.

    package main

    import (
        "fmt"
        "github.com/hurngchunlee/gotoon"
    )

    func main() {

        // step 1: initialize Toon with authentication credentials.
        toon := gotoon.Toon{
            TenantID:       "eneco",
            Username:       "myEnecoUsername",
            Password:       "myEnecoPassword",
            ConsumerKey:    "ToonAPIConsumerKey",
            ConsumerSecret: "ToonAPIConsumerSecret",
        }

        // step 2: retrieve the agreements.
        agreements, _ := toon.GetAgreements()

        // step 3: retrieve status and information.
        for _, agreement := range agreements {
            status, _ := toon.GetStatus(agreement)
            fmt.Printf("%+v\n", status)
        }

    }

*/
package gotoon
