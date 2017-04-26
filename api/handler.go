package api

import (
	"net/http"

	"encoding/json"

	"fmt"

	"github.com/concourse/pagerhub/cmd"
	"github.com/concourse/pagerhub/pagerduty"
	"github.com/tedsuo/rata"
)

func NewHandler(opts *cmd.Opts, pagerdutyClient *pagerduty.Client) (http.Handler, error) {
	handlers := map[string]http.Handler{
		HealthCheck: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
		}),
		Webhook: &GithubSignatureMiddleware{
			Inner: &WebhookHandler{
				PagerdutyIntegrationKey: opts.PagerdutyIntegrationKey,
				PagerdutyClient:         pagerdutyClient,
			},
			GithubWebhookSecret: opts.GithubWebhookSecret,
		},
	}

	return rata.NewRouter(Routes, handlers)
}

type WebhookHandler struct {
	PagerdutyIntegrationKey string
	PagerdutyClient         *pagerduty.Client
}

func (w *WebhookHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	const githubEventHeader = "X-GitHub-Event"
	eventType := r.Header.Get(githubEventHeader)

	if GithubWebhookEvent(eventType) != GithubWebhookEventIssues {
		rw.WriteHeader(http.StatusAccepted)
		rw.Write([]byte("ignoring: event was not of type 'issue'"))
		return
	}

	defer r.Body.Close()

	var event GithubIssuesEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("could not decode event into JSON: " + err.Error()))
		return
	}

	if event.Action != GithubIssuesEventActionOpened {
		rw.WriteHeader(http.StatusAccepted)
		rw.Write([]byte("ignoring: issue was not 'opened'"))
		return
	}

	e := pagerduty.Event{
		Action:     pagerduty.ActionTrigger,
		RoutingKey: w.PagerdutyIntegrationKey,
		Payload: pagerduty.Payload{
			Summary:  event.Issue.User.Login + " reported an issue",
			Source:   event.Issue.HTMLURL,
			Severity: pagerduty.SeverityWarning,
		},
	}

	err = w.PagerdutyClient.Enqueue(e)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("could not trigger pagerduty alert: " + err.Error() + "\n"))
		rw.Write([]byte("pagerduty event: " + fmt.Sprintf("%v", e)))
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

type GithubWebhookEvent string

const (
	GithubWebhookEventIssues GithubWebhookEvent = "issues"
)

type GithubIssuesEvent struct {
	Action GithubIssuesEventAction `json:"action"`
	Issue  GithubIssue             `json:"issue"`
}

type GithubIssuesEventAction string

const (
	GithubIssuesEventActionOpened GithubIssuesEventAction = "opened"
)

type GithubIssue struct {
	ID      int        `json:"id"`
	Number  int        `json:"number"`
	Title   string     `json:"title"`
	URL     string     `json:"url"`
	HTMLURL string     `json:"html_url"`
	Body    string     `json:"body"`
	User    GithubUser `json:"user"`
}

type GithubUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}
