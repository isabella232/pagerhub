package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/pkg/errors"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Enqueue(e Event) error {
	b := new(bytes.Buffer)

	err := json.NewEncoder(b).Encode(e)
	if err != nil {
		return errors.Wrap(err, "encoding pagerduty event failed")
	}

	req, err := http.NewRequest("POST", "https://events.pagerduty.com/v2/enqueue", b)
	if err != nil {
		return errors.Wrap(err, "generating request to pagerduty failed")
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "making request to pagerduty failed")
	}

	if res.StatusCode >= http.StatusBadRequest {
		r, err := httputil.DumpResponse(res, true)
		if err != nil {
			return errors.Wrap(err, "dumping response failed")
		}

		return fmt.Errorf("bad pagerduty request (%i): %s", res.StatusCode, r)
	}

	return nil
}

type Event struct {
	RoutingKey string  `json:"routing_key" required:"true"`
	Action     Action  `json:"event_action" required:"true"`
	Payload    Payload `json:"payload" required:"true"`
}

type Payload struct {
	Summary  string   `json:"summary" required:"true"`
	Source   string   `json:"source" required:"true"`
	Severity Severity `json:"severity" required:"true"`
}

type Action string

const (
	ActionTrigger     Action = "trigger"
	ActionAcknowledge Action = "acknowledge"
	ActionResolve     Action = "resolve"
)

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityError    Severity = "error"
	SeverityWarning  Severity = "warning"
	SeverityInfo     Severity = "info"
)
