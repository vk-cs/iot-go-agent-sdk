package iot_go_agent_sdk

// Status represents current status of agent or device
type Status string

const (
	Online  Status = "online"
	Offline Status = "offline"
)
