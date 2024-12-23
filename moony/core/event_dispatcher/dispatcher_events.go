package event_dispatcher

type EventType string

const (
	OnServerStarted EventType = "OnServerStarted"
	OnServerStopped EventType = "OnServerStopped"
)
