package iot_go_agent_sdk

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/mailru/surgemq/message"
	"github.com/mailru/surgemq/service"

	"github.com/vk-cs/iot-go-agent-sdk/mqtt"
)

const (
	EventTopic               = "iot/event/fmt/json"
	LogTopic                 = "iot/log/fmt/json"
	CommandTopic             = "iot/cmd/agent/%d/fmt/json"
	CommandAgentStatusTopic  = "iot/cmd/agent/%d/status/fmt/json"
	CommandDeviceStatusTopic = "iot/cmd/device/%d/status/fmt/json"
)

type MQTTClient struct {
	connector *service.Client
	agentID   int64
	login     string
	password  string
	host      string
}

func (c *MQTTClient) Connect() error {
	msg := message.NewConnectMessage()
	if err := msg.SetWillQos(message.QosAtLeastOnce); err != nil {
		return fmt.Errorf("failed to set will QoS: %w", err)
	}
	if err := msg.SetVersion(4); err != nil {
		return fmt.Errorf("failed to set version: %w", err)
	}

	if err := msg.SetClientId([]byte(uuid.New().String())); err != nil {
		return fmt.Errorf("failed to set client ID: %w", err)
	}
	msg.SetWillFlag(true)
	msg.SetUsername([]byte(c.login))
	msg.SetPassword([]byte(c.password))
	msg.SetCleanSession(true)

	return c.connector.Connect(c.host, msg)
}

func (c *MQTTClient) Disconnect() {
	c.connector.Disconnect()
}

func (c *MQTTClient) PublishEvent(packetID uint16, eventMessage *mqtt.EventMessage) error {
	return c.publishMessage(packetID, EventTopic, eventMessage)
}

func (c *MQTTClient) PublishLog(packetID uint16, logMessage *mqtt.LogMessage) error {
	return c.publishMessage(packetID, LogTopic, logMessage)
}

func (c *MQTTClient) PublishCommand(packetID uint16, agentID int64, command *mqtt.CommandMessage) error {
	topic := fmt.Sprintf(CommandTopic, agentID)
	return c.publishMessage(packetID, topic, command)
}

func (c *MQTTClient) PublishAgentCommandStatus(packetID uint16, command *mqtt.CommandStatusMessage) error {
	topic := fmt.Sprintf(CommandAgentStatusTopic, c.agentID)
	return c.publishMessage(packetID, topic, command)
}

func (c *MQTTClient) PublishDeviceCommandStatus(
	packetID uint16, deviceID int64, command *mqtt.CommandStatusMessage,
) error {
	topic := fmt.Sprintf(CommandDeviceStatusTopic, deviceID)
	return c.publishMessage(packetID, topic, command)
}

func (c *MQTTClient) Subscribe(subscriber service.Subscriber) error {
	topic := fmt.Sprintf(CommandTopic, c.agentID)
	msg := message.NewSubscribeMessage()
	if err := msg.AddTopic([]byte(topic), message.QosAtLeastOnce); err != nil {
		return fmt.Errorf("failed to set topic: %w", err)
	}
	if err := c.connector.Subscribe(msg, subscriber); err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}
	return nil
}

func (c *MQTTClient) publishMessage(packetID uint16, topic string, payload mqtt.Message) error {
	msg := message.NewPublishMessage()
	msg.SetPacketId(packetID)
	msg.SetRetain(true)

	if err := msg.SetQoS(message.QosAtLeastOnce); err != nil {
		return fmt.Errorf("failed to set will QoS: %w", err)
	}

	if err := msg.SetTopic([]byte(topic)); err != nil {
		return fmt.Errorf("failed to set topic: %w", err)
	}

	if err := payload.ValidateToPublish(); err != nil {
		return fmt.Errorf("failed to validate message: %w", err)
	}

	marshaledPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	msg.SetPayload(marshaledPayload)

	return c.connector.Publish(msg, nil)
}

func NewMQTTClient(agentID int64, login, password, host string) (*MQTTClient, error) {
	newClient := &MQTTClient{
		connector: &service.Client{},
		agentID:   agentID,
		login:     login,
		password:  password,
		host:      host,
	}
	if err := newClient.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to broker: %w", err)
	}
	return newClient, nil
}
