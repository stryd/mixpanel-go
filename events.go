package mixpanel

import (
	"bytes"
	"crypto/md5"
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
		Location *EventLocation

		// Time the event occurred. Set to nil to use the current time.
		Time *time.Time

		// CustomProperties that wished to be tracked as fields of the event
		CustomProperties map[string]interface{}
	}
	EventLocation struct {
		City string
		// Region usually refers to the state. Can be any string?
		Region string
		// Country upper case 2 character ISO country code
		Country string
	}
)

type eventsClient struct {
	c internalClient
	// TODO batch import requests
}

// Track create an event to current distinct id
func (ec *eventsClient) Track(e Event) error {
	params := encodeToWireFormat(e, ec.c.config.Token)
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
		if len(e.InsertID) == 0 {
			// import endpoint requires insert ids to be included
			f := md5.New()
			_, _ = f.Write([]byte(fmt.Sprintf("%v", e)))
			e.InsertID = fmt.Sprintf("%x", f.Sum(nil))
		}
		payloadEvents[i] = encodeToWireFormat(e, "")
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

func encodeToWireFormat(e Event, token string) map[string]interface{} {
	props := map[string]interface{}{
		"distinct_id": e.DistinctID,
	}
	if len(token) > 0 {
		props["token"] = token
	}
	if e.IPV4Address != nil {
		props["ip"] = e.IPV4Address
	}
	if e.Location != nil {
		if len(e.Location.City) > 0 {
			props["$city"] = e.Location.City
		}
		if len(e.Location.Region) > 0 {
			props["$region"] = e.Location.Region
		}
		if len(e.Location.Country) == 2 {
			props["mp_country_code"] = strings.ToUpper(e.Location.Country)
		}
	}
	if e.Time != nil {
		props["time"] = e.Time.Unix()
	}
	if len(e.InsertID) > 0 {
		props["$insert_id"] = e.InsertID
	}

	for key, value := range e.CustomProperties {
		props[key] = value
	}

	params := map[string]interface{}{
		"event":      e.Name,
		"properties": props,
	}
	return params
}
