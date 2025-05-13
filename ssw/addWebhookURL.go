package ssw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func AddWebhookURL(APIKey, currentDevice, webhookURL string) (err error) {
	// Send the HTTP POST request
	body := []byte(`{
    "device_name":"` + currentDevice + `",
    "url":"` + webhookURL + `"
}`)
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		"https://app.supersimplewhats.com/v1/devices/hook_endpoints/add",
		bytes.NewBuffer(body),
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

	response := httpResponse{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("error adding webhook URL: %s", response.Code)
	}

	return nil
}
