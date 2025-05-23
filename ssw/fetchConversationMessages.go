package ssw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type ConversationData struct {
	ConversationID  string        `json:"conversation_id"`
	DestinationName string        `json:"destination_name"`
	DestinationType string        `json:"destination_type"`
	LastMessageDate string        `json:"last_message_date"`
	MessageCount    int           `json:"message_count"`
	Messages        []Message     `json:"messages"`
	Participants    []Participant `json:"participants"`
}

type Message struct {
	ID          string    `json:"id"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	DeletedAt   *string   `json:"deleted_at"` // Using pointer to handle null
	MessageDate time.Time `json:"message_date"`
	ContactName string    `json:"contact_name"`
	FromMe      bool      `json:"from_me"`
	Message     string    `json:"message"`
	Type        string    `json:"type"`
}

func (m Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	return json.Marshal(&struct {
		MessageDate string `json:"message_date"`
		*Alias
	}{
		MessageDate: m.MessageDate.Format("2006-01-02T15:04:05Z"),
		Alias:       (*Alias)(&m),
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface for Message
func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	aux := &struct {
		MessageDate string `json:"message_date"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.MessageDate != "" {
		parsedTime, err := time.Parse("2006-01-02T15:04:05Z", aux.MessageDate)
		if err != nil {
			return err
		}
		m.MessageDate = parsedTime
	}

	return nil
}

type Participant struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func FetchConversationMessages(APIKey, currentRecipient string) (output []Message, err error) {
	// Send the HTTP GET request
	body := []byte(`{"conversation_id": "` + currentRecipient + `", "limit": 10}`)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://app.supersimplewhats.com/v1/conversations/list_messages", bytes.NewBuffer(body))
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
		return nil, fmt.Errorf("error listing messages: %s", response.Code)
	}

	// decode response.Data into ConversationData
	var conversationData ConversationData
	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling conversation data: %v", err)
	}

	err = json.Unmarshal([]byte(jsonData), &conversationData)
	if err != nil {
		return nil, fmt.Errorf("error decoding conversation data: %v", err)
	}

	output = conversationData.Messages

	// sort output by message date
	sort.Slice(output, func(i, j int) bool {
		return output[i].MessageDate.Before(output[j].MessageDate)
	})

	return output, nil
}
