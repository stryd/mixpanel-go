package events

type EventName string

// holds all valid mixpanel event names
var EventNames = struct {
	TestServerEvent EventName
	RunUploaded     EventName
}{"test_server_event", "run_uploaded"}
