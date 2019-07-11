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
	app := app.App{}

	settings, err := app.Setup(plgver + " (" + version + ")")
	if err != nil {
		log.Printf("Unable to start the app: %s", err)
		os.Exit(1)
	}

	app.Run(settings)
}
