package mixpanel

import (
	"testing"
	"time"
)

func TestMock(t *testing.T) {
	m := NewMockClient()

	future := time.Date(2015, time.October, 21, 0, 0, 0, 0, time.UTC)
	ip := "1.2.3.4"

	_ = m.Track(Event{
		Name:        "Sign Up",
		DistinctID:  "1",
		Time:        &future,
		IPV4Address: &ip,
		CustomProperties: map[string]interface{}{
			"Referred By": "Friend",
			"URL":         "mixpanel.com/signup",
		},
	})

	_ = m.Import([]Event{{
		Name:        "Sign Up",
		DistinctID:  "2",
		Time:        &future,
		IPV4Address: &ip,
		CustomProperties: map[string]interface{}{
			"Referred By": "Friend",
			"URL":         "mixpanel.com/signup",
		},
	}, {
		Name:       "Purchase",
		DistinctID: "2",
		Time:       &future,
		InsertID:   "935d87b1-00cd-41b7-be34-b9d98dd08b42",
		CustomProperties: map[string]interface{}{
			"Item":   "Coffee",
			"Amount": 5.0,
		},
	}})

	if len(m.eventsByID) != 2 {
		t.Fatalf("incorrect number of distinct ids saved")
	}
	if len(m.eventsByID["1"]) != 1 {
		t.Fatalf("incorrect number of events for user 1")
	}
	if len(m.eventsByID["2"]) != 2 {
		t.Fatalf("incorrect number of events for user 2")
	}

	// Output:
	// {
	//	"1": [
	//		{
	//			"Name": "Sign Up",
	//			"DistinctID": "1",
	//			"InsertID": "",
	//			"IPV4Address": "1.2.3.4",
	//			"Location": null,
	//			"Time": "2015-10-21T00:00:00Z",
	//			"CustomProperties": {
	//				"Referred By": "Friend",
	//				"URL": "mixpanel.com/signup"
	//			}
	//		}
	//	],
	//	"2": [
	//		{
	//			"Name": "Sign Up",
	//			"DistinctID": "2",
	//			"InsertID": "2c2098df8333d0ee01b6d7e747094d85",
	//			"IPV4Address": "1.2.3.4",
	//			"Location": null,
	//			"Time": "2015-10-21T00:00:00Z",
	//			"CustomProperties": {
	//				"Referred By": "Friend",
	//				"URL": "mixpanel.com/signup"
	//			}
	//		},
	//		{
	//			"Name": "Purchase",
	//			"DistinctID": "2",
	//			"InsertID": "935d87b1-00cd-41b7-be34-b9d98dd08b42",
	//			"IPV4Address": null,
	//			"Location": null,
	//			"Time": "2015-10-21T00:00:00Z",
	//			"CustomProperties": {
	//				"Amount": 5,
	//				"Item": "Coffee"
	//			}
	//		}
	//	]
	// }
}
