package main

import (
	"log"
	"os"

	"jbrodriguez/unbalance/server/src/app"
)

// Version -
var Version string

func main() {
	app := app.App{}

	settings, err := app.Setup(Version)
	if err != nil {
		log.Printf("Unable to start the app: %s", err)
		os.Exit(1)
	}

	app.Run(settings)
}
