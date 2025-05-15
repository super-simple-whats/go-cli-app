# SuperSimpleWhats API - example Golang cli app

This is a simple example of a Golang CLI application that uses the SuperSimpleWhats API to send WhatsApp messages.

You are not supposed to use this code in production, but it is a good starting point to understand how easy is to use the SuperSimpleWhats API.

## Prerequisites

- Go 1.18 or later
- A [SuperSimpleWhats](https://supersimplewhats.com/) account
- A WhatsApp Business account already running in a device
- A [NGROK](https://ngrok.com/) account (for webhook testing)

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/super-simple-whats/go-cli-app.git
    cd go-cli-app
    ```

2. Install the required dependencies:
    ```bash
    go mod tidy
    ```

3. Create a `.env` file in the root directory and add your SuperSimpleWhats API key:
    ```bash
    API_KEY=<your_api_key>
    HOOKS_HOST=localhost:8080
    HOOKS_PATH=/hooks
    ```

4. Install the `ngrok` to expose your local server to the internet.
    ##### For Mac
    ```bash
    brew install ngrok
    ```
    ##### For Linux
    ```bash
    curl -sSL https://ngrok-agent.s3.amazonaws.com/ngrok.asc \
      | sudo tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null \
      && echo "deb https://ngrok-agent.s3.amazonaws.com buster main" \
      | sudo tee /etc/apt/sources.list.d/ngrok.list \
      && sudo apt update \
      && sudo apt install ngrok
    ```
    ##### For Windows
    If you are using Windows, you can download the installer from the [ngrok website](https://ngrok.com/download).

5. Start ngrok to expose your local server, for example:
    ```bash
    ngrok http 8080
    ```

6. 🥳🎉 That's it! Now run the application and it will look for your registered devices and webhook endpoints. Follow the instructions and you are all set!
    ```bash
    go run .
    ```

