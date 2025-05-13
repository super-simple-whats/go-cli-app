package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"msg_repo_cli/ssw"

	"github.com/skip2/go-qrcode"
)

var (
	webhookReady = make(chan bool)
	qrcodeReady  = make(chan bool)
)

func main() {
	loadEnv()
	// Start the webhook server in a goroutine
	go startWebhookServer()

	// Wait for the webhook server to be ready
	<-webhookReady
	fmt.Println("Webhook server is ready.")

	var err error
	devices, err = ssw.FetchDevices(APIKey)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Main input loop
	scanner := bufio.NewScanner(os.Stdin)

	if len(devices) == 0 {
		registerDevice()
	}

	if currentDevice == "" && len(devices) > 0 {
		chooseDevice()
	}

	URLs, err := ssw.FetchWebhookURLs(APIKey, currentDevice)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if len(URLs) == 0 {
		addWebhookURL()
	}

	mainMenu()

	for {
		// fmt.Printf("[%s] You: ", currentDevice)
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()

		// Check for commands
		if strings.HasPrefix(input, "/") {
			parts := strings.SplitN(input, " ", 2)
			command := parts[0]

			switch command {
			case "/menu":
				mainMenu()
			case "/devices":
				chooseDevice()
			case "/list-conversations":
				chooseConversation()
			case "/exit":
				fmt.Println("Exiting...")
				return

			default:
				fmt.Printf("Unknown command: %s\n", command)
			}

			continue
		}

		// Send regular message
		if input != "" {
			if err := ssw.SendMessage(APIKey, currentDevice, currentRecipient, input); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
}

func mainMenu() {
	fmt.Println("Golang CLI WhatsApp Application")
	fmt.Println("Type your message and press Enter to send.")
	fmt.Println("Use /menu to show this message again.")
	fmt.Println("Use /devices to list all devices.")
	fmt.Println("Use /add-webhook-url to add a webhook URL.")
	fmt.Println("Use /register-device [device name] to register a new device.")
	fmt.Println("Use /list-conversations to list all conversations.")
	fmt.Println("Type /exit to quit.")
	fmt.Println("--------------------------------------------------")
	fmt.Printf("[%s] You: ", currentDevice)
}

func registerDevice() {
	var err error
	var code string

	clearScreen()
	fmt.Println("No devices found with the given API Key. Registering a new device...")
	fmt.Println("--------------------------------------------------")

	fmt.Print("Enter device name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	deviceName := scanner.Text()

	code, err = ssw.RegisterDevice(APIKey, deviceName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Device registered successfully. Please use the generated code to connect your WhatsApp: %s\n", code)
	qrCode, err := qrcode.New(code, qrcode.Medium)
	if err != nil {
		fmt.Printf("Failed to generate QR code: %v\n", err)
	} else {
		fmt.Println(qrCode.ToSmallString(false))
		fmt.Println("After scanning the QR code, please restart the application.")
	}

	<-qrcodeReady
}

func chooseDevice() {
	clearScreen()

	fmt.Println("Available Devices:")
	for i, device := range devices {
		fmt.Printf(" %d. %s\n", i+1, device)
	}

	fmt.Print("\nEnter device number (1-", len(devices), "): ")
	var choice int
	_, err := fmt.Scanf("%d", &choice)
	if err != nil || choice < 1 || choice > len(devices) {
		fmt.Println("Invalid selection. Please try again.")
		return
	}

	// Set the selected device (assuming you have a variable for this)
	currentDevice = devices[choice-1]
	fmt.Printf("Selected device: %s\n", currentDevice)

	// You might want to return the selected device or set it to a global variable
	// depending on how your program is structured
	time.Sleep(2 * time.Second)
	clearScreen()
}

func addWebhookURL() {
	var err error
	var webhookURL string

	clearScreen()
	fmt.Println("No webhook URLs found. Adding a new webhook URL...")
	fmt.Println("--------------------------------------------------")

	fmt.Print("Enter webhook URL: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	webhookURL = scanner.Text()

	err = ssw.AddWebhookURL(APIKey, currentDevice, webhookURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Webhook URL added successfully.")
	time.Sleep(2 * time.Second)
	clearScreen()
	mainMenu()
}

func chooseConversation() {
	var err error
	var recipients []ssw.Conversation

	recipients, err = ssw.FetchConversations(currentDevice, APIKey)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	clearScreen()

	fmt.Println("Available Conversations:")
	for i, recipient := range recipients {
		formattedDateStr := recipient.LastMessageDate
		formattedDate, err := time.Parse("2006-01-02T15:04:05Z", recipient.LastMessageDate)
		if err == nil {
			formattedDateStr = formattedDate.Format("2006-01-02 15:04:05")
		}
		fmt.Printf(" %d. %s [%s] - %s\n", i+1, recipient.ChatName, recipient.ChatType, formattedDateStr)
	}

	fmt.Print("\nEnter conversation number (1-", len(recipients), "): ")
	var choice int
	_, err = fmt.Scanf("%d", &choice)
	if err != nil || choice < 1 || choice > len(recipients) {
		fmt.Println("Invalid selection. Please try again.")
		return
	}

	currentRecipient = recipients[choice-1].ID
	fmt.Printf("Selected conversation: %s\n", currentRecipient)

	time.Sleep(1 * time.Second)
	clearScreen()

	messages, err = ssw.FetchConversationMessages(APIKey, currentRecipient)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	showConversation()
}

func showConversation() {
	// Implement the logic to show the conversation with the given recipientID
	clearScreen()

	for _, message := range messages {
		formattedDateStr := message.MessageDate
		formattedDate, err := time.Parse("2006-01-02T15:04:05Z", message.MessageDate)
		if err == nil {
			formattedDateStr = formattedDate.Format("2006-01-02 15:04:05")
		}

		name := message.ContactName
		if message.FromMe {
			name = "You"
		}
		fmt.Printf("[%s] %s: %s\n", formattedDateStr, name, message.Message)
	}

	fmt.Printf("[%s] You: ", currentDevice)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
