# mixpanel-go

Mixpanel Go Client and Stryd Mixpanel event names

## Usage

```go
import (
  "github.com/stryd/mixpanel-go"
  "github.com/stryd/mixpanel-go/events"
)
```

## Examples

Initialize the tracking client:

```go
mixpanelClient := mixpanel.Init(<MIXPANEL_API_TOKEN>)
```

Track an event:

```go
eventParams := EventProperties{
  "run_id": "123456789"
}

err := mixpanelClient.Track(<USER_ID>, events.EventNames.SomeEventName, &eventParams)
```
