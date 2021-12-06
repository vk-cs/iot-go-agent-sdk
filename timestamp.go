package iot_go_agent_sdk

import "time"

// Now returns current timestamp in microseconds (default time format in platform API)
func Now() int64 {
	now := time.Now().UnixNano() / 1000
	return now
}
