package dispatcher

// this file is outdated and is not used right now
// because all plugins events are kinda generated in real time

type EventType string

const (
	OnServerStarted EventType = "OnServerStarted"
	OnServerStopped EventType = "OnServerStopped"
)
