package ssw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Conversation struct {
	ID              string `json:"id"`
	ChatName        string `json:"chat_name"`
	ChatType        string `json:"chat_type"`
	LastMessageDate string `json:"last_message_date"`
}

func FetchConversations(currentDevice, APIKey string) (output []Conversation, err error) {
	// Send the HTTP GET request
	client := &http.Client{}
	body := []byte(`{"platforms": ["whatsapp"], "device_keys": ["` + currentDevice + `"], "start_date": "2025-01-01T00:00:00Z", "end_date": "2025-12-31T23:59:59Z"}`)
	fmt.Printf("\n\nBody: %v\n", string(body))

	req, err := http.NewRequest("POST", "https://app.supersimplewhats.com/v1/conversations/list", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", APIKey)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status code: %d", resp.StatusCode)
	}

	response := httpResponse{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("error listing conversations: %s", response.Code)
	}
	if response.Data == nil {
		return nil, fmt.Errorf("error listing conversations: empty data")
	}

	// Convert the data to a slice of conversations
	data, ok := response.Data.([]any)
	if !ok {
		return nil, fmt.Errorf("error parsing conversations data")
	}
	for _, conv := range data {
		convMap, ok := conv.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("error parsing conversation data")
		}
		convStruct := Conversation{
			ID:              convMap["id"].(string),
			ChatName:        convMap["chat_name"].(string),
			ChatType:        convMap["chat_type"].(string),
			LastMessageDate: convMap["last_message_date"].(string),
		}
		output = append(output, convStruct)
	}

	if len(output) == 0 {
		return nil, fmt.Errorf("error listing conversations: empty conversations")
	}

	return output, nil
}
