package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/concourse/pagerhub/api"
	"github.com/jessevdk/go-flags"
	"github.com/vito/twentythousandtonnesofcrudeoil"
	"github.com/concourse/pagerhub/cmd"
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

	handler, err := api.NewHandler(opts)
	if err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(opts.Port), handler))
}
