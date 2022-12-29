package domain

type MQTTMsg struct {
	Type   string `json:"type"`
	Action string `json:"action"`
}
