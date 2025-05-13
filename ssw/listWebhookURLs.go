package ssw

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Hook struct {
	ID         string `json:"id"`
	DeviceName string `json:"device_name"`
	URL        string `json:"url"`
}

func FetchWebhookURLs(APIKey, currentDevice string) (URLs []Hook, err error) {
	url := "https://app.supersimplewhats.com/v1/devices/hook_endpoints/list"
	method := "POST"

	payload := strings.NewReader(`{
    "device_name":"` + currentDevice + `"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", APIKey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	response := httpResponse{}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return
	}

	if !response.Success {
		return
	}

	URLs = make([]Hook, 0)

	jsonData, err := json.Marshal(response.Data)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonData, &URLs)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}

	return
}
