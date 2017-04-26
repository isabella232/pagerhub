package main

import (
	"log"
	"net/http"
	"os"
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
		os.Exit(1)
	}

	p := pagerduty.NewClient()

	handler, err := api.NewHandler(opts, p)
	if err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(opts.Port), handler))
}
