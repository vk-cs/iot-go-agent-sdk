package mqtt

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEvent(t *testing.T) {
	eventTagInt := &EventTag{ID: 10, Value: 100, Timestamp: 1}
	eventTagStr := &EventTag{ID: 11, Value: "abc", Timestamp: 1}
	eventMessage := &EventMessage{Tags: []*EventTag{eventTagInt, eventTagStr}}
	assert.NoError(t, eventMessage.ValidateToPublish())
}

func TestValidateCommand(t *testing.T) {
	tagInt := &Tag{ID: 10, Value: "foobar"}
	tagStr := &Tag{ID: 20, Value: 0}
	tagBool := &Tag{ID: 30, Value: true}

	commandInt := &Command{ID: "CommandInt", Timestamp: 1, Tags: []*Tag{tagInt}}
	commandStr := &Command{ID: "CommandStr", Timestamp: 1, Tags: []*Tag{tagStr}}
	commandStrBool := &Command{ID: "CommandStrBool", Timestamp: 1, Tags: []*Tag{tagStr, tagBool}}

	deviceInt := &DeviceCommand{ID: 100, Command: commandInt}
	deviceStr := &DeviceCommand{ID: 200, Command: commandStr}

	assert.Error(t, (&CommandMessage{}).ValidateToPublish())
	assert.Error(t, (&CommandMessage{Devices: []*DeviceCommand{deviceStr}, Command: commandInt}).ValidateToPublish())
	assert.Error(t, (&CommandMessage{Devices: []*DeviceCommand{}}).ValidateToPublish())
	assert.Error(t, (&CommandMessage{Devices: []*DeviceCommand{deviceInt, deviceStr}}).ValidateToPublish())

	assert.NoError(t, (&CommandMessage{Command: commandStrBool}).ValidateToPublish())
	assert.NoError(t, (&CommandMessage{Devices: []*DeviceCommand{deviceInt}}).ValidateToPublish())
}

func TestValidateCommandStatus(t *testing.T) {
	commandDone := &CommandStatusMessage{ID: "CommandInt", Status: Done, Timestamp: 1}
	commandFailed := &CommandStatusMessage{ID: "CommandInt", Status: Failed, Reason: "test", Timestamp: 1}
	assert.NoError(t, commandDone.ValidateToPublish())
	assert.NoError(t, commandFailed.ValidateToPublish())

	commandDoneWithReason := &CommandStatusMessage{ID: "CommandInt", Status: Done, Reason: "test", Timestamp: 1}
	commandFailedWithoutReason := &CommandStatusMessage{ID: "CommandInt", Status: Failed, Timestamp: 1}
	assert.Error(t, commandDoneWithReason.ValidateToPublish())
	assert.Error(t, commandFailedWithoutReason.ValidateToPublish())
}

func TestDumpLogMessage(t *testing.T) {
	logMessage := &LogMessage{Message: "foo", Timestamp: 1}
	body, err := json.Marshal(logMessage)
	assert.NoError(t, err)

	expected := make(map[string]interface{})
	expected["msg"] = "foo"
	expected["timestamp"] = 1
	expectedJSON, _ := json.Marshal(expected)

	assert.Equalf(t, body, expectedJSON, "marshal log message is incorrect")
}
