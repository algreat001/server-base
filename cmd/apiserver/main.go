package main

import (
	"flag"
	"log"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver"
	//"server/internal/app/apiserver"
)

func main() {
	flag.Parse()

	if err := apiserver.Start(false); err != nil {
		log.Fatal(err)
	}
}
