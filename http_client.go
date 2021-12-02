package iot_go_agent_sdk

import (
	"github.com/go-openapi/strfmt"

	"github.com/vk-cs/iot-go-agent-sdk/gen/swagger/http_client/client"
)

const (
	DefaultHTTPHost = "api-iot.mcs.mail.ru"
)

func NewHTTPClient() *client.HTTP {
	transportCfg := client.DefaultTransportConfig()
	transportCfg.WithHost(DefaultHTTPHost)
	transportCfg.WithSchemes([]string{"https"})
	return client.NewHTTPClientWithConfig(strfmt.Default, transportCfg)
}