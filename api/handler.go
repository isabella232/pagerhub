package api

import (
	"net/http"

	"github.com/concourse/pagerhub/cmd"
	"github.com/tedsuo/rata"
)

func NewHandler(opts *cmd.Opts) (http.Handler, error) {
	handlers := map[string]http.Handler{
		HealthCheck: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
		}),
		Webhook: &GithubSignatureMiddleware{
			GithubWebhookSecret: opts.Github.WebhookSecretToken,
			Inner: &WebhookHandler{},
		},
	}

	return rata.NewRouter(Routes, handlers)
}

type WebhookHandler struct {
}

func (w *WebhookHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusAccepted)
}
