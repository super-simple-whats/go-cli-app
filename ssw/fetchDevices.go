package ssw

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchDevices(APIKey string) (devices []string, err error) {
	// Send the HTTP GET request
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://app.supersimplewhats.com/v1/devices/list", nil)
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
		return nil, fmt.Errorf("error listing devices: %s", response.Code)
	}

	for _, device := range response.Data.([]any) {
		// fmt.Printf("\nDevice: %v\n", device)
		if deviceName, ok := device.(map[string]any)["name"].(string); ok {
			devices = append(devices, deviceName)
		} else {
			return nil, fmt.Errorf("error parsing device name")
		}
	}

	return devices, nil
}
