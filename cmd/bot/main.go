package main

import (
	"log"
	"os"

	"github.com/antonbatrakov/chatgpt-telegram-bot/pkg/chatgpt"
	"github.com/antonbatrakov/chatgpt-telegram-bot/pkg/telegram"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram struct {
		Token        string  `yaml:"token"`
		AllowedUsers []int64 `yaml:"allowed_users"`
	} `yaml:"telegram"`
	ChatGPT struct {
		APIKey      string  `yaml:"api_key"`
		Model       string  `yaml:"model"`
		MaxTokens   int     `yaml:"max_tokens"`
		Temperature float64 `yaml:"temperature"`
	} `yaml:"chatgpt"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}

func main() {
	// Load configuration
	configFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Initialize ChatGPT client
	chatGPTClient := chatgpt.NewClient(
		config.ChatGPT.APIKey,
		config.ChatGPT.Model,
		config.ChatGPT.MaxTokens,
		config.ChatGPT.Temperature,
	)

	// Initialize Telegram bot
	bot, err := telegram.NewBot(config.Telegram.Token, chatGPTClient, config.Telegram.AllowedUsers)
	if err != nil {
		log.Fatalf("Error creating Telegram bot: %v", err)
	}

	// Start the bot
	log.Println("Starting bot...")
	bot.Start()
}
