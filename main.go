package main

import (
	"github.com/lithictech/runtime-js-env/jsenv"
	"log"
)

func main() {
	err := jsenv.InstallAt("index.html", jsenv.DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}
}
