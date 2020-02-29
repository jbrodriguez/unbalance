package main

import (
	"log"
	"os"

	"unbalance/app"
)

// Version -
var version string
var plgver string

func main() {
	appv := app.App{}

	settings, err := appv.Setup(plgver + " (" + version + ")")
	if err != nil {
		log.Printf("Unable to start the app: %s", err)
		os.Exit(1)
	}

	appv.Run(settings)
}
