package main

import (
	"log"

	"github.com/simonnik/GB_observability/hw3/app"
)

func main() {
	a := app.App{}
	if closer, err := a.Init(); err != nil {
		log.Fatal(err)
	} else {
		defer closer.Close()
	}

	if err := a.Serve(); err != nil {
		log.Fatal(err)
	}
}
