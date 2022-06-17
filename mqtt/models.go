package mqtt

type Command struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Tags      []*Tag `json:"tags"`
}

type DeviceCommand struct {
	ID      int64    `json:"device_id"`
	Command *Command `json:"command"`
}

type Tag struct {
	ID    int64       `json:"id"`
	Value interface{} `json:"value"`
}

type EventTag struct {
	ID        int64       `json:"id"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
}

type CommandStatus string

const (
	Skipped  CommandStatus = "skipped"
	Received               = "received"
	Failed                 = "failed"
	Done                   = "done"
)
