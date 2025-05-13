package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"msg_repo_cli/ssw"
)

// Start webhook server to listen for incoming messages
func startWebhookServer() {
	http.HandleFunc(hooksPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var event ssw.WebhookEvent
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&event); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if event.DeviceName != currentDevice {
			return
		}

		if event.EventName == "message_received" {
			fmt.Printf("\n\n Webhook event: %v\n", event)
			jsonData, err := json.Marshal(event.Data)
			if err != nil {
			}
			var message ssw.Message
			if err := json.Unmarshal(jsonData, &message); err != nil {
			}
			messages = append(messages, message)
			showConversation()
		} else if event.EventName == "message_sent" {
			fmt.Printf("\n\n Webhook event: %v\n", event)
			jsonData, err := json.Marshal(event.Data)
			if err != nil {
			}
			var message ssw.Message
			if err := json.Unmarshal(jsonData, &message); err != nil {
			}
			messages = append(messages, message)
			showConversation()
		} else if event.EventName == "qr_code" {
			// Cast the data to string and print as QR code
			// if qrData, ok := event.Data.(string); ok {
			// 	qrCode, err := qrcode.New(qrData, qrcode.Medium)
			// 	if err != nil {
			// 		fmt.Printf("Failed to generate QR code: %v\n", err)
			// 	} else {
			// 		fmt.Println(qrCode.ToSmallString(false))
			// 	}
			// } else {
			// 	fmt.Printf("QR code data is not a string: %v\n", event.Data)
			// }
			fmt.Printf("Event: %v\n", event)
			qrcodeReady <- true
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Println("Webhook server listening on http://" + hooksHost + hooksPath)
	webhookReady <- true
	if err := http.ListenAndServe(hooksHost, nil); err != nil {
		log.Fatalf("Failed to start webhook server: %v", err)
	}
}
