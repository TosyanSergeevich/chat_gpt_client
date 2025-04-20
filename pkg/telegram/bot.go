package telegram

import (
	"fmt"
	"sync"

	"github.com/antonbatrakov/chatgpt-telegram-bot/pkg/chatgpt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	chatGPT      *chatgpt.Client
	allowedUsers map[int64]bool
	sessions     map[int64][]chatgpt.Message
	mu           sync.Mutex
}

func NewBot(token string, chatGPT *chatgpt.Client, allowedUsers []int64) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("error creating bot: %v", err)
	}

	allowedUsersMap := make(map[int64]bool)
	for _, userID := range allowedUsers {
		allowedUsersMap[userID] = true
	}

	return &Bot{
		bot:          bot,
		chatGPT:      chatGPT,
		allowedUsers: allowedUsersMap,
		sessions:     make(map[int64][]chatgpt.Message),
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Check if user is allowed
		if len(b.allowedUsers) > 0 && !b.allowedUsers[update.Message.From.ID] {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, you are not authorized to use this bot.")
			b.bot.Send(msg)
			continue
		}

		// Handle commands
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				b.handleStart(update.Message)
			case "reset":
				b.handleReset(update.Message)
			}
			continue
		}

		// Handle text messages
		if update.Message.Text != "" {
			go b.handleTextMessage(update.Message)
		}

		// Handle photo messages
		if update.Message.Photo != nil {
			go b.handlePhotoMessage(update.Message)
		}
	}
}

func (b *Bot) handleStart(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! I'm your ChatGPT assistant. Send me a message or use /reset to start a new chat.")
	msg.ReplyMarkup = b.createKeyboard()
	b.bot.Send(msg)
}

func (b *Bot) handleReset(message *tgbotapi.Message) {
	b.mu.Lock()
	b.sessions[message.Chat.ID] = nil
	b.mu.Unlock()

	msg := tgbotapi.NewMessage(message.Chat.ID, "Chat session reset. Starting a new conversation.")
	msg.ReplyMarkup = b.createKeyboard()
	b.bot.Send(msg)
}

func (b *Bot) createKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/reset"),
		),
	)
}

func (b *Bot) handleTextMessage(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Processing your message...")
	sentMsg, _ := b.bot.Send(msg)

	b.mu.Lock()
	session := b.sessions[message.Chat.ID]
	if session == nil {
		session = []chatgpt.Message{}
	}

	// Add user message to session
	session = append(session, chatgpt.Message{
		Role:    "user",
		Content: message.Text,
	})

	response, err := b.chatGPT.SendMessage(session)

	if err == nil {
		// Add assistant response to session
		session = append(session, chatgpt.Message{
			Role:    "assistant",
			Content: response,
		})
		b.sessions[message.Chat.ID] = session
	}
	b.mu.Unlock()

	if err != nil {
		editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, sentMsg.MessageID, fmt.Sprintf("Error: %v", err))
		b.bot.Send(editMsg)
		return
	}

	editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, sentMsg.MessageID, response)
	b.bot.Send(editMsg)
}

func (b *Bot) handlePhotoMessage(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Processing your image...")
	sentMsg, _ := b.bot.Send(msg)

	// Get the largest photo size
	photo := message.Photo[len(message.Photo)-1]
	file, err := b.bot.GetFile(tgbotapi.FileConfig{FileID: photo.FileID})
	if err != nil {
		editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, sentMsg.MessageID, fmt.Sprintf("Error getting file: %v", err))
		b.bot.Send(editMsg)
		return
	}

	// Download the photo
	photoURL := file.Link(b.bot.Token)
	response, err := b.chatGPT.SendImageMessage([]chatgpt.ImageMessage{
		{
			Role: "user",
			Content: []chatgpt.ContentObject{
				{
					Type: "text",
					Text: message.Caption,
				},
				{
					Type: "image_url",
					ImageURL: chatgpt.ImageURL{
						URL: photoURL,
					},
				},
			},
		},
	})

	if err != nil {
		editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, sentMsg.MessageID, fmt.Sprintf("Error: %v", err))
		b.bot.Send(editMsg)
		return
	}

	editMsg := tgbotapi.NewEditMessageText(message.Chat.ID, sentMsg.MessageID, response)
	b.bot.Send(editMsg)
}
