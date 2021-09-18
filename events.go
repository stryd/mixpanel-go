package mixpanel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type (
	EventsAPI interface {
		// Import a batch of events to Mixpanel (usually for server side tracking)
		Import(events []Event) error

		// Track events to Mixpanel (usually for client-side tracking)
		Track(e Event) error
	}
	// Event A mixpanel event
	Event struct {
		// Name of the event being tracked.
		Name string

		// DistinctID of the entity that produced this event. (usually a user)
		DistinctID string

		// InsertID is a unique identifier for the event, used for duplication. Leave blank to have one generated.
		InsertID string

		// IPV4Address of the user when they created the event. Only provide one of IPV4Address or Location.
		IPV4Address *string
		// Location is an optional field to specify where this event occurred. Only provide one of IPV4Address or Location.
		Location *struct {
			City    string
			Region  string
			Country string
		}

		// Time the event occurred. Set to nil to use the current time.
		Time *time.Time

		// CustomProperties that wished to be tracked as fields of the event
		CustomProperties map[string]interface{}
	}
)

type eventsClient struct {
	c internalClient
	// TODO batch import requests
}

// Track create an event to current distinct id
func (ec *eventsClient) Track(e Event) error {
	props := map[string]interface{}{
		"token":       ec.c.config.Token,
		"distinct_id": e.DistinctID,
	}
	if e.IPV4Address != nil {
		props["ip"] = e.IPV4Address
	}
	if e.Time != nil {
		props["time"] = e.Time.Unix()
	}

	for key, value := range e.CustomProperties {
		props[key] = value
	}

	params := map[string]interface{}{
		"event":      e.Name,
		"properties": props,
	}
	dataBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}

	values := url.Values{}
	values.Set("data", string(dataBytes))
	// TODO support ip and verbose

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/track", ec.c.config.ApiUrl), strings.NewReader(values.Encode()))
	req.Header.Add(Headers.Accept, MimeTypes.TextPlain)
	req.Header.Add(Headers.ContentType, MimeTypes.XWWWFormURLEncoded)

	return ec.c.send(req)
}

func (ec *eventsClient) Import(events []Event) error {
	payloadEvents := make([]map[string]interface{}, len(events))
	for i, e := range events {
		props := map[string]interface{}{
			"distinct_id": e.DistinctID,
		}
		if e.IPV4Address != nil {
			props["ip"] = e.IPV4Address
		}
		if e.Time != nil {
			props["time"] = e.Time.Unix()
		}
		for key, value := range e.CustomProperties {
			props[key] = value
		}

		payloadEvents[i] = map[string]interface{}{
			"event":      e.Name,
			"properties": props,
		}
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payloadEvents); err != nil {
		return err
	}

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/import?strict=1&project_id=%s", ec.c.config.ApiUrl, ec.c.config.ProjectID), &buf)
	req.Header.Add(Headers.Accept, MimeTypes.ApplicationJSON)
	req.Header.Add(Headers.ContentType, MimeTypes.ApplicationJSON)

	return ec.c.send(req)
}
