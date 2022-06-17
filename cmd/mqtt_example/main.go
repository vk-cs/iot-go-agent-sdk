package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/mailru/surgemq/message"
	"github.com/vk-cs/iot-go-agent-sdk/mqtt"

	sdk "github.com/vk-cs/iot-go-agent-sdk"
	"github.com/vk-cs/iot-go-agent-sdk/gen/swagger/http_client/client/agents"
)

const DefaultMQTTHost = "tcp://mqtt-api-iot.mcs.mail.ru:1883"

func MQTTExample(login, password, host string) error {
	fmt.Println("Starting new MQTT example func")

	fmt.Println("Getting latest config")
	httpClient := sdk.NewHTTPClient()
	agentConfig, err := httpClient.Agents.GetAgentConfig(
		&agents.GetAgentConfigParams{Version: nil, Context: context.Background()},
		httptransport.BasicAuth(login, password),
	)
	if err != nil {
		return fmt.Errorf("failed to get agent config through API: %w", err)
	}
	agentID := *agentConfig.Payload.Agent.ID

	fmt.Println("Creating MQTT agent")
	mqttClient := sdk.NewMQTTClient(agentID, host)

	if err := mqttClient.Connect(login, password); err != nil {
		return fmt.Errorf("example error during MQTT client connection: %w", err)
	}
	defer mqttClient.Disconnect()

	fmt.Println("Sending event (status of agent)")
	statusTag, _ := sdk.FindTagByPath(agentConfig.Payload.Agent.Tag.Children, sdk.StatusTagPath)
	bootstrappingEventTag := &mqtt.EventTag{
		ID:        *statusTag.ID,
		Value:     "bootstrapping",
		Timestamp: sdk.Now(),
	}
	bootstrappingEventMessage := &mqtt.EventMessage{Tags: []*mqtt.EventTag{bootstrappingEventTag}}
	if err := mqttClient.PublishEvent(bootstrappingEventMessage); err != nil {
		return fmt.Errorf("example error during first event publishing: %w", err)
	}

	fmt.Println("Marking agent as online")
	now := sdk.Now()
	updatedTag, _ := sdk.FindTagByPath(agentConfig.Payload.Agent.Tag.Children, sdk.ConfigUpdatedAtTagPath)
	onlineEventMessage := &mqtt.EventMessage{Tags: []*mqtt.EventTag{
		{ID: *statusTag.ID, Value: "online", Timestamp: now},
		{ID: *updatedTag.ID, Value: now, Timestamp: now},
	}}
	if err := mqttClient.PublishEvent(onlineEventMessage); err != nil {
		return fmt.Errorf("example error during online event publishing: %w", err)
	}

	fmt.Println("Subscribing for agent commands")
	if err := mqttClient.Subscribe(&subscriberExample{MQTTClient: mqttClient}); err != nil {
		return fmt.Errorf("example error during subscribing: %w", err)
	}

	fmt.Println("Sending temperature from thermometer of first device")
	rand.Seed(time.Now().UnixNano())
	temperature := rand.Intn(30)
	temperatureTag, found := sdk.FindTagByPath(
		agentConfig.Payload.Agent.Devices[0].Tag.Children,
		[]string{"temperature"},
	)
	if found == false {
		return fmt.Errorf("no temperature tag for tesing")
	}
	temperatureEventMessage := &mqtt.EventMessage{Tags: []*mqtt.EventTag{
		{ID: *temperatureTag.ID, Value: temperature, Timestamp: sdk.Now()},
	}}
	if err := mqttClient.PublishEvent(temperatureEventMessage); err != nil {
		return fmt.Errorf("example error during temperature event publishing: %w", err)
	}

	fmt.Println("Sending log message")
	logMessage := &mqtt.LogMessage{
		Message:   fmt.Sprintf("temperature is %d", temperature),
		Timestamp: sdk.Now(),
	}
	if err := mqttClient.PublishLog(logMessage); err != nil {
		return fmt.Errorf("example error during log publishing: %w", err)
	}

	fmt.Println("Finishing MQTT example func")
	return nil
}

type subscriberExample struct {
	MQTTClient *sdk.MQTTClient
}

func (s *subscriberExample) OnPublish(msg *message.PublishMessage) error {
	fmt.Println("Starting OnPublish handler")
	commandMessage := mqtt.CommandMessage{}
	if err := json.Unmarshal(msg.Payload(), &commandMessage); err != nil {
		return fmt.Errorf("failed unmarshaling command message: %w", err)
	}
	if commandMessage.Devices != nil {
		for _, deviceCommand := range commandMessage.Devices {
			fmt.Println("Notifying about device command received")
			if err := s.MQTTClient.PublishDeviceCommandStatus(
				deviceCommand.ID,
				&mqtt.CommandStatusMessage{
					ID:        deviceCommand.Command.ID,
					Status:    mqtt.Received,
					Timestamp: sdk.Now(),
				},
			); err != nil {
				return fmt.Errorf("failed sending device command status recieved: %w", err)
			}

			// do some work ...

			fmt.Println("Notifying about device command done")
			if err := s.MQTTClient.PublishDeviceCommandStatus(
				deviceCommand.ID,
				&mqtt.CommandStatusMessage{
					ID:        deviceCommand.Command.ID,
					Status:    mqtt.Done,
					Timestamp: sdk.Now(),
				},
			); err != nil {
				return fmt.Errorf("failed sending device command status recieved: %w", err)
			}
		}

		if commandMessage.Command != nil {
			fmt.Println("Notifying about agent command received")
			if err := s.MQTTClient.PublishAgentCommandStatus(
				&mqtt.CommandStatusMessage{
					ID:        commandMessage.Command.ID,
					Status:    mqtt.Received,
					Timestamp: sdk.Now(),
				},
			); err != nil {
				return fmt.Errorf("failed sending agent command status recieved: %w", err)
			}

			// do some work ...

			fmt.Println("Notifying about agent command failed (for test reason)")
			if err := s.MQTTClient.PublishAgentCommandStatus(
				&mqtt.CommandStatusMessage{
					ID:        commandMessage.Command.ID,
					Status:    mqtt.Failed,
					Reason:    "Failed for test reason",
					Timestamp: sdk.Now(),
				},
			); err != nil {
				return fmt.Errorf("failed sending agent command status as failed: %w", err)
			}
		}
	}
	fmt.Println("Finishing OnPublish handler")
	return nil
}

func (s *subscriberExample) OnComplete(msg, ack message.Message, err error) error {
	return nil
}

func main() {
	login := flag.String("login", "", "login for agent")
	password := flag.String("password", "", "password for agent")
	host := flag.String("host", DefaultMQTTHost, "mqtt host with `tcp://{host}:{port}` format")

	flag.Parse()
	if *login == "" {
		fmt.Println("Missing login param")
		os.Exit(1)
	}
	if *password == "" {
		fmt.Println("Missing password param")
		os.Exit(1)
	}

	if err := MQTTExample(*login, *password, *host); err != nil {
		fmt.Println("Example error: %w", err)
		os.Exit(1)
	}
	fmt.Println("Example completed successfully")
}
