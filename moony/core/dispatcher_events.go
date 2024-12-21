package core

type EventType string

const (
	OnServerStarted EventType = "OnServerStarted"
	OnServerStopped EventType = "OnServerStopped"
)
