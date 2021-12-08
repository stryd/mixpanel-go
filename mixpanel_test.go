package mixpanel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type testServer struct {
	server      *httptest.Server
	lastRequest *http.Request
	lastEvents  []Event
	lastToken   string
}

var (
	token     = "api_token"
	account   = "service_account"
	secret    = "service_secret"
	projectID = "test_project"
)

func setup() (*testServer, API) {
	ts := testServer{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.lastRequest = r
		switch r.Header.Get(Headers.ContentType) {
		case MimeTypes.XWWWFormURLEncoded:
			_ = r.ParseForm()
			data := r.FormValue("data")
			var m map[string]interface{}
			if err := json.Unmarshal([]byte(data), &m); err != nil {
				w.WriteHeader(400)
				_, _ = w.Write([]byte("0"))
				return
			}
			ts.lastEvents = []Event{ts.decodeEventFromWire(m)}
		case MimeTypes.ApplicationJSON:
			var events []map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
				w.WriteHeader(400)
				_, _ = w.Write([]byte("0"))
				return
			}
			for i := range events {
				ts.lastEvents = append(ts.lastEvents, ts.decodeEventFromWire(events[i]))
			}
		}

		w.WriteHeader(200)
		_, _ = w.Write([]byte("1"))
	}))
	ts.server = server
	return &ts, New(WithToken(token), WithApiUrl(ts.server.URL), WithServiceAccount(account, secret), WithProjectID(projectID))
}

func (server *testServer) decodeEventFromWire(m map[string]interface{}) Event {
	props := m["properties"].(map[string]interface{})
	e := Event{
		Name:             m["event"].(string),
		CustomProperties: make(map[string]interface{}),
	}
	for k, v := range props {
		switch k {
		case "distinct_id":
			e.DistinctID = v.(string)
		case "insert_id":
			e.InsertID = v.(string)
		case "token":
			server.lastToken = v.(string)
		default:
			e.CustomProperties[k] = v
		}
	}
	return e
}

func TestTrack(t *testing.T) {
	ts, m := setup()
	defer ts.server.Close()

	event := Event{
		Name:       "Signed Up",
		DistinctID: "1337",
		CustomProperties: map[string]interface{}{
			"Referred By": "Friend",
		},
	}
	if err := m.Track(event); err != nil {
		t.Error(err)
	}

	if ts.lastRequest.Method != http.MethodPost {
		t.Errorf("incorrect http method.\n expected %s: actual %s", http.MethodPost, ts.lastRequest.Method)
	}
	if len(ts.lastEvents) != 1 {
		t.Errorf("incorrect number of events recieved.\n expected 1: actual %d", len(ts.lastEvents))
	}
	if !reflect.DeepEqual(event, ts.lastEvents[0]) {
		t.Errorf("recieved event does not match sent event. %+v != %+v", event, ts.lastEvents[0])
	}
	if !strings.EqualFold(ts.lastRequest.URL.Path, "/track") {
		t.Errorf("incorrect endpoint called.\n expected %s: actual %s", "/track", ts.lastRequest.URL.Path)
	}
	if !strings.EqualFold(token, ts.lastToken) {
		t.Errorf("token not sent in track call.\n expected %s: actual %s", token, ts.lastToken)
	}
}

func TestImport(t *testing.T) {
	ts, m := setup()
	defer ts.server.Close()

	events := []Event{{
		Name:       "Signed Up",
		DistinctID: "1337",
		InsertID:   "123",
		CustomProperties: map[string]interface{}{
			"Referred By": "Friend",
		},
	}, {
		Name:       "Purchase",
		DistinctID: "1337",
		InsertID:   "456",
		CustomProperties: map[string]interface{}{
			"Amount": 5.0,
		},
	}}
	if err := m.Import(events); err != nil {
		t.Error(err)
	}

	if ts.lastRequest.Method != http.MethodPost {
		t.Errorf("incorrect http method.\n expected %s: actual %s", http.MethodPost, ts.lastRequest.Method)
	}
	if len(ts.lastEvents) != len(events) {
		t.Errorf("incorrect number of events recieved.\n expected %d: actual %d", len(events), len(ts.lastEvents))
	}
	for i := range events {
		for j := range ts.lastEvents {
			if events[i].InsertID == ts.lastEvents[i].InsertID {
				if !reflect.DeepEqual(events[i], ts.lastEvents[j]) {
					t.Errorf("recieved event does not match sent event. %+v != %+v", events[i], ts.lastEvents[j])
				}
			}
		}
	}
	if !strings.EqualFold(ts.lastRequest.URL.Path, "/import") {
		t.Errorf("incorrect endpoint called.\n expected %s: actual %s", "/import", ts.lastRequest.URL.Path)
	}
	expectedQueryParams := fmt.Sprintf("strict=1&project_id=%s", projectID)
	if !strings.EqualFold(expectedQueryParams, ts.lastRequest.URL.RawQuery) {
		t.Errorf("incorrect query params passed.\n expected %s: actual %s", expectedQueryParams, ts.lastRequest.URL.RawQuery)
	}
	u, p, ok := ts.lastRequest.BasicAuth()
	if !ok {
		t.Errorf("auth secret not included")
	}
	if !strings.EqualFold(account, u) {
		t.Errorf("incorrect auth account.\n expected %s: actual %s", account, u)
	}
	if !strings.EqualFold(secret, p) {
		t.Errorf("incorrect auth secret.\n expected %s: actual %s", secret, p)
	}
}
