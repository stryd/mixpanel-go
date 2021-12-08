package mixpanel

import (
	"fmt"
	"time"
)

func ExampleMock() {
	m := NewMock()

	t, _ := time.Parse(time.RFC3339, "2016-03-03T15:17:53+01:00")

	_ = m.Update("1", &Update{
		Operation: "$set",
		Timestamp: &t,
		IP:        "127.0.0.1",
		Properties: map[string]interface{}{
			"custom_field": "cool!",
		},
	})

	_ = m.Track("1", "Sign Up", &Event{
		IP: "1.2.3.4",
		Properties: map[string]interface{}{
			"from": "email",
		},
	})

	fmt.Println(m)

	// Output:
	// 1:
	//   ip: 127.0.0.1
	//   time: 2016-03-03T15:17:53+01:00
	//   properties:
	//     custom_field: cool!
	//   events:
	//     Sign Up:
	//       IP: 1.2.3.4
	//       Timestamp:
	//       from: email
}
