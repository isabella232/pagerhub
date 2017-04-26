package cmd

type Opts struct {
	Port                    int    `short:"p" long:"port" default:"8080" description:"Port to run webserver on"`
	GithubWebhookSecret     string `long:"github-webhook-secret" description:"Github Webhook Secret" required:"true"`
	PagerdutyIntegrationKey string `long:"pagerduty-integration-key" description:"Pagerduty Integration Key" required:"true"`
}
