package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/concourse/pagerhub/api"
	"github.com/concourse/pagerhub/cmd"
	"github.com/concourse/pagerhub/pagerduty"
	"github.com/jessevdk/go-flags"
	"github.com/vito/twentythousandtonnesofcrudeoil"
)

func main() {
	opts := &cmd.Opts{}

	parser := flags.NewParser(opts, flags.Default)
	parser.NamespaceDelimiter = "-"

	twentythousandtonnesofcrudeoil.TheEnvironmentIsPerfectlySafe(parser, "")

	_, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	p := pagerduty.NewClient()

	handler, err := api.NewHandler(opts, p)
	if err != nil {
		log.Fatal(err)
	}

	port := strconv.Itoa(opts.Port)
	log.Println("starting server on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
