package domain

type MQTTMsg struct {
	Type   string `json:"type"`
	Action string `json:"action"`
}

type UploadImageMsg struct {
	Type      string `json:"type"`
	ImageData string `json:"image_data"`
}
