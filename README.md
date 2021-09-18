# mixpanel-go

Mixpanel Go Client

## Usage

```go
import (
  "github.com/stryd/mixpanel-go"
  "github.com/stryd/mixpanel-go/events"
)
```

## Examples

Initialize the client:

```go
mixpanelClient := mixpanel.New("api-token", "")
```

Track an event:

```go
event := mixpanel.Event{
	Properties: map[string]interface{}{
        "run_id": "123456789",
    }
}

err := mixpanelClient.Track("distinct-id", "your-event", &event)
```
