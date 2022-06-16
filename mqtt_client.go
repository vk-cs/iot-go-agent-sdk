package iot_go_agent_sdk

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

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
	lock            *sync.RWMutex
	connector       *service.Client
	agentID         int64
	host            string
	randomGenerator *rand.Rand
}

func (c *MQTTClient) Connect(login, password string) error {
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
	msg.SetUsername([]byte(login))
	msg.SetPassword([]byte(password))
	msg.SetCleanSession(true)

	return c.connector.Connect(c.host, msg)
}

func (c *MQTTClient) Disconnect() {
	c.connector.Disconnect()
}

func (c *MQTTClient) PublishEvent(eventMessage *mqtt.EventMessage) error {
	return c.publishMessage(EventTopic, eventMessage)
}

func (c *MQTTClient) PublishLog(logMessage *mqtt.LogMessage) error {
	return c.publishMessage(LogTopic, logMessage)
}

func (c *MQTTClient) PublishCommand(agentID int64, command *mqtt.CommandMessage) error {
	topic := fmt.Sprintf(CommandTopic, agentID)
	return c.publishMessage(topic, command)
}

func (c *MQTTClient) PublishAgentCommandStatus(command *mqtt.CommandStatusMessage) error {
	topic := fmt.Sprintf(CommandAgentStatusTopic, c.agentID)
	return c.publishMessage(topic, command)
}

func (c *MQTTClient) PublishDeviceCommandStatus(deviceID int64, command *mqtt.CommandStatusMessage) error {
	topic := fmt.Sprintf(CommandDeviceStatusTopic, deviceID)
	return c.publishMessage(topic, command)
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

func (c *MQTTClient) publishMessage(topic string, payload mqtt.Message) error {
	msg := message.NewPublishMessage()
	msg.SetPacketId(c.makePacketID())
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

func (c *MQTTClient) makePacketID() uint16 {
	c.lock.Lock()
	defer c.lock.Unlock()

	return uint16(c.randomGenerator.Intn(math.MaxUint16))
}

func NewMQTTClient(agentID int64, host string) *MQTTClient {
	return &MQTTClient{
		lock:            &sync.RWMutex{},
		connector:       &service.Client{},
		agentID:         agentID,
		host:            host,
		randomGenerator: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
