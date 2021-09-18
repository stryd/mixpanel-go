# mixpanel-go

Mixpanel Go Client

## Usage

`go get github.com/stryd/mixpanel-go`

```go
import (
  "github.com/stryd/mixpanel-go"
)
```
## Examples

### Initializing a new client

```go
mixpanelAPI := mixpanel.New(mixpanel.WithToken("api-token"))
```

There are multiple options that can be provided in any order
```go
mixpanel.WithToken("api_token")
mixpanel.WithApiUrl("https://new-api.mixpanel.com")
mixpanel.WithSecret("service_account_secret")
mixpanel.WithProjectID("my-project")
```

### Tracking a single event

This is usually used client-side but can be used from servers as well. The event must have occurred with in the last 5 days if a time is provided.

```go
event := mixpanel.Event{
	DistinctID: "1234",
	Name: "your-event",
	Properties: map[string]interface{}{
		"run_id": "123456789",
	}
}

if err := mixpanelClient.Track(event); err != nil {
	log.Fatal(err)
}
```

### Importing multiple events

Should only be used server side. Can track up to 2000 events as long as those 2000 events don't exceed the size limit
of a single import request (2MB uncompressed). Insert IDs will be calculated if not provided for each event.

```go
now := time.Now()

events := []mixpanel.Event{{
	DistinctID: "1234",
	Name: "run-completed",
	Time: &now,
	Properties: map[string]interface{}{
		"run_id": "123456789",
	}
},{
	DistinctID: "1234", 
	Name: "run-deleted",
	Time: &now,
	Properties: map[string]interface{}{
		"run_id": "123456789",
	}
}}

if err := mixpanelClient.Import(events); err != nil {
	log.Fatal(err)
}
```
