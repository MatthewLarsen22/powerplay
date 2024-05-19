package chat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jak103/powerplay/internal/server/apis"
	"github.com/jak103/powerplay/internal/server/services/auth"
)

// Most simple message structure
type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

// Pseudo database of messages
var messages []Message

func init() {
	apis.RegisterHandler(fiber.MethodGet, "/chat", auth.Public, getMessages)
	apis.RegisterHandler(fiber.MethodPost, "/chat", auth.Public, postMessage)
}

// Return all messages belonging to a certain recipient
func getMessages(context *fiber.Ctx) error {
	name := context.Query("name")
	var chatHistory []Message

	for _, msg := range messages {
		if msg.To == name || msg.From == name {
			chatHistory = append(chatHistory, msg)
		}
	}

	return context.JSON(chatHistory)
}

// Simple post message function which adds to our messages array.
func postMessage(context *fiber.Ctx) error {
	var message Message

	err := context.BodyParser(&message)
	if err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":    "Cannot parse JSON",
			"received": message.Content,
		})
	}

	isMessageValid := message.Content != "" &&
		message.From != "" &&
		message.To != ""

	if !isMessageValid {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad info.  'To', 'From', and 'Content' must not be empty...",
		})
	}

	messages = append(messages, message)

	return context.SendStatus(fiber.StatusCreated)
}
