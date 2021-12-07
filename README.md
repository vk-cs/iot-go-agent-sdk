# iot-go-agent-sdk
SDK for building agents for VK CS IoT Platform based on Go programming language

# Usage

```go
package main

import (
	"context"

	httptransport "github.com/go-openapi/runtime/client"
	
	sdk "github.com/vk-cs/iot-go-agent-sdk"
	"github.com/vk-cs/iot-go-agent-sdk/gen/swagger/http_client/client/agents"
	"github.com/vk-cs/iot-go-agent-sdk/gen/swagger/http_client/client/events"
	"github.com/vk-cs/iot-go-agent-sdk/gen/swagger/http_client/models"
)

func main() {
	cli := sdk.NewHTTPClient()

	// get agent config
	resp, _ := cli.Agents.GetAgentConfig(&agents.GetAgentConfigParams{
		// get latest version
		Version: nil,
		Context: context.Background(),
	}, httptransport.BasicAuth("login", "password"))
	cfg := resp.Payload

	// get status tag directly from config
	statusTag, _ := sdk.FindTagByPath(cfg.Agent.Tag, sdk.StatusTagPath)
	
	// or use wrapper for quick access
	tree := sdk.NewTagTree(cfg.Agent.Tag)
	statusTagNode, _ := tree.GetStatusTag()
	statusTag = statusTagNode.Tag
	cfgVersion, _ := tree.GetConfigVersionTag()
	cfgUpdatedAt, _ := tree.ConfigUpdatedAtTagPath()
	
	// send to platform new agent status and cfg version
	now := sdk.Now()
	cli.Events.AddEvent(events.AddEventParams{
		Context: context.Background(),
		Body: &models.AddEvent{
			Tags: []*models.TagValueObject{
				// set agent status
				{
					ID: statusTag.ID,
					Timestamp: &now,
					Value: sdk.Online,
				},
				// set agent config version
				{
					ID: cfgVersion.Tag.ID,
					Timestamp: &now,
					Value: cfg.Version,
				},
				{
					ID: cfgUpdatedAt.Tag.ID,
					Timestamp: &now,
					Value: &now,
				},
			},
		},
	}, httptransport.BasicAuth("login", "password"))
	
	// same for device
	deviceStatusTag, _ := sdk.FindTagByPath(cfg.Agent.Devices[0].Tag, sdk.StatusTagPath)

	cli.Events.AddEvent(events.AddEventParams{
		Context: context.Background(),
		Body: &models.AddEvent{
			Tags: []*models.TagValueObject{
				{
					ID: deviceStatusTag.ID,
					Timestamp: &now,
					Value: sdk.Online,
				},				
			},
		},
	}, httptransport.BasicAuth("login", "password"))	
}
```

# Development

Regenerate http client from swagger api spec
```bash
make generate
```
