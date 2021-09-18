package mixpanel

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
)

// MockClient Mocked version of Mixpanel intended for unit tests only
type MockClient struct {
	eventsByID map[string][]Event
}

func NewMockClient() *MockClient {
	return &MockClient{
		eventsByID: make(map[string][]Event),
	}
}

func (m *MockClient) Track(e Event) error {
	m.eventsByID[e.DistinctID] = append(m.eventsByID[e.DistinctID], e)
	return nil
}

func (m *MockClient) Import(es []Event) error {
	for _, e := range es {
		if len(e.InsertID) == 0 {
			f := md5.New()
			_, _ = f.Write([]byte(fmt.Sprintf("%v", e)))
			e.InsertID = fmt.Sprintf("%x", f.Sum(nil))
		}
		m.eventsByID[e.DistinctID] = append(m.eventsByID[e.DistinctID], e)
	}
	return nil
}

func (m *MockClient) String() string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "\t")
	if err := enc.Encode(m.eventsByID); err != nil {
		return "ERROR"
	}
	return buf.String()
}
