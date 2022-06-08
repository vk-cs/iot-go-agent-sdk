package mqtt

import "errors"

type Message interface {
	ValidateToPublish() error
}

type EventMessage struct {
	Tags []*EventTag `json:"tags"`
}

func (m *EventMessage) ValidateToPublish() error {
	return nil
}

type CommandMessage struct {
	Command *Command         `json:"command"`
	Devices []*DeviceCommand `json:"devices"`
}

func (m *CommandMessage) ValidateToPublish() error {
	if m.Command == nil && m.Devices == nil {
		return errors.New("empty command message")
	}
	if m.Devices != nil && len(m.Devices) != 1 {
		return errors.New("you can publish only one device command per message")
	}
	if m.Command != nil && m.Devices != nil {
		return errors.New("sending more than one command per request is prohibited")
	}
	return nil
}

type CommandStatusMessage struct {
	ID        string        `json:"id"`
	Status    CommandStatus `json:"status"`
	Reason    string        `json:"reason"`
	Timestamp int64         `json:"timestamp"`
}

func (m *CommandStatusMessage) ValidateToPublish() error {
	if m.Status == Failed && m.Reason == "" {
		return errors.New("failed command status without set reason")
	}
	if m.Status != Failed && m.Reason != "" {
		return errors.New("reason have been set for not failed status")
	}
	return nil
}

type LogMessage struct {
	Message   string `json:"msg"`
	Timestamp int64  `json:"timestamp"`
}

func (m *LogMessage) ValidateToPublish() error {
	return nil
}
