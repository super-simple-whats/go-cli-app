package ssw

import (
	"bytes"
	"fmt"
	"net/http"
)

// Outgoing message structure
type MessageRequest struct {
	DeviceKey string `json:"device_key"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

// Send a message to the specified recipient
func SendMessage(APIKey, currentDevice, recipient, message string) error {
	// Create the request payload

	// Send the HTTP POST request
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		"https://app.supersimplewhats.com/v1/messages/send/"+currentDevice+"/"+recipient,
		bytes.NewBuffer([]byte(message)),
	)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", APIKey)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status code: %d", resp.StatusCode)
	}

	return nil
}
