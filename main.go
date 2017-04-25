package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/concourse/pagerhub/api"
	"github.com/jessevdk/go-flags"
	"github.com/vito/twentythousandtonnesofcrudeoil"
)

func main() {
	opts := &Opts{}

	parser := flags.NewParser(opts, flags.Default)
	parser.NamespaceDelimiter = "-"

	twentythousandtonnesofcrudeoil.TheEnvironmentIsPerfectlySafe(parser, "")

	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	handler, err := api.NewHandler()
	if err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(opts.Port), handler))
}

type Opts struct {
	Port   int `short:"p" long:"port" default:"8080" description:"Port to run webserver on"`
	Github struct {
		WebhookSecretToken    string `long:"webhook-secret" description:"Webhook secret" required:"true"`
		IntegrationPrivateKey string `long:"integration-private-key" description:"Integration private key" required:"true"`
		OAuth                 struct {
			ClientID     string `long:"client-id" description:"OAuth Client ID" required:"true"`
			ClientSecret string `long:"client-secret" description:"OAuth Client Secret" required:"true"`
		} `group:"OAuth" namespace:"oauth"`
	} `group:"Github" namespace:"github"`
}
