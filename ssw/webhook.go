package ssw

// Incoming webhook event structure
type WebhookEvent struct {
	DeviceName string `json:"device_name"`
	EventName  string `json:"event_name"`
	Data       any    `json:"data"`
}
