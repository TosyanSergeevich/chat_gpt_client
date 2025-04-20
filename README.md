# ChatGPT Telegram Bot

A Telegram bot that allows you to chat with ChatGPT and send images for analysis.

## Features

- Chat with ChatGPT through Telegram
- Send images to ChatGPT for analysis
- User access control
- Configurable settings

## Prerequisites

- Go 1.16 or higher
- Telegram Bot Token (from [@BotFather](https://t.me/BotFather))
- OpenAI API Key

## Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/chatgpt-telegram-bot.git
cd chatgpt-telegram-bot
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure the bot:
   - Copy `config/config.yaml.example` to `config/config.yaml`
   - Fill in your Telegram bot token and OpenAI API key
   - Add allowed user IDs if you want to restrict access

4. Build and run the bot:
```bash
go build -o bot cmd/bot/main.go
./bot
```

## Configuration

Edit `config/config.yaml` to customize the bot's behavior:

```yaml
telegram:
  token: "YOUR_TELEGRAM_BOT_TOKEN"
  allowed_users: []  # List of allowed user IDs

chatgpt:
  api_key: "YOUR_CHATGPT_API_KEY"
  model: "gpt-4-vision-preview"  # Using vision model for image support
  max_tokens: 1000
  temperature: 0.7

server:
  port: 8080
```

## Usage

1. Start a chat with your bot on Telegram
2. Send text messages to chat with ChatGPT
3. Send images with optional captions for image analysis

## Security

- The bot supports user access control through the `allowed_users` configuration
- Keep your API keys secure and never commit them to version control
- Consider using environment variables for sensitive information

## License

MIT 