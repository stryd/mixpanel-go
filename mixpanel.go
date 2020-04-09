// Fork of https://github.com/dukex/mixpanel
package mixpanel

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/stryd/mixpanel-go/events"
)

type MixpanelError struct {
	URL string
	Err error
}

func (err *MixpanelError) Cause() error {
	return err.Err
}

func (err *MixpanelError) Error() string {
	return "mixpanel: " + err.Err.Error()
}

type ErrTrackFailed struct {
	Body string
	Resp *http.Response
}

// The Mixpanel struct stores the mixpanel endpoint and the project token
type mixpanel struct {
	Client *http.Client
	Token  string
	ApiURL string
}

type EventProperties map[string]interface{}

// An update of a user in mixpanel
type Update struct {
	// Update operation such as "$set", "$update" etc.
	Operation string

	// Custom properties. At least one must be specified.
	Properties EventProperties
}

func (err *ErrTrackFailed) Error() string {
	return fmt.Sprintf("Mixpanel did not return 1 when tracking: %s", err.Body)
}

type Mixpanel interface {
	// Create a mixpanel event
	Track(userId string, eventName events.EventName, props *EventProperties) error

	// Set properties for a mixpanel user.
	Update(userId string, u *Update) error

	Alias(userId string, newId string) error
}

// When alias is called, Mixpanel will create a mapping between the Mixpanel generated distinct_id, and your unique identifier.
// Only call Alias once; an alias can only point to one Mixpanel distinct_id.
// Subsequently, when identify is called, you pass your identifier and Mixpanel will connect it with the original distinct_id.
// see https://help.mixpanel.com/hc/en-us/articles/115004497803-Identity-Management-Best-Practices for more on identity management best practices
func (m *mixpanel) Alias(userId, newId string) error {
	props := map[string]interface{}{
		"token":       m.Token,
		"distinct_id": userId,
		"alias":       newId,
	}

	params := map[string]interface{}{
		"event":      "$create_alias",
		"properties": props,
	}

	return m.send("track", params, false)
}

// Track sends an event to mixpanel tied to a specific user ID
func (m *mixpanel) Track(userId string, eventName events.EventName, props *EventProperties) error {
	eventInfo := map[string]interface{}{
		"token":       m.Token,
		"distinct_id": userId,
	}

	for key, value := range *props {
		eventInfo[key] = value
	}

	params := map[string]interface{}{
		"event":      eventName,
		"properties": eventInfo,
	}

	autoGeolocate := true

	return m.send("track", params, autoGeolocate)
}

// Updates a user in mixpanel. See
// https://mixpanel.com/help/reference/http#people-analytics-updates
func (m *mixpanel) Update(userId string, u *Update) error {
	params := map[string]interface{}{
		"$token":       m.Token,
		"$distinct_id": userId,
	}

	params[u.Operation] = u.Properties

	autoGeolocate := true

	return m.send("engage", params, autoGeolocate)
}

func (m *mixpanel) to64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func (m *mixpanel) send(eventType string, params interface{}, autoGeolocate bool) error {
	data, err := json.Marshal(params)

	if err != nil {
		return err
	}

	url := m.ApiURL + "/" + eventType + "?data=" + m.to64(data)

	if autoGeolocate {
		url += "&ip=1"
	}

	wrapErr := func(err error) error {
		return &MixpanelError{URL: url, Err: err}
	}

	resp, err := m.Client.Get(url)

	if err != nil {
		return wrapErr(err)
	}

	defer resp.Body.Close()

	body, bodyErr := ioutil.ReadAll(resp.Body)

	if bodyErr != nil {
		return wrapErr(bodyErr)
	}

	if strBody := string(body); strBody != "1" && strBody != "1\n" {
		return wrapErr(&ErrTrackFailed{Body: strBody, Resp: resp})
	}

	return nil
}

// Init returns the client instance.
func Init(token string) Mixpanel {
	return NewFromClient(http.DefaultClient, token)
}

// Creates a client instance using the specified client instance.
func NewFromClient(c *http.Client, token string) Mixpanel {
	return &mixpanel{
		Client: c,
		Token:  token,
		ApiURL: "https://api.mixpanel.com",
	}
}
