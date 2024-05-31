package chat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jak103/powerplay/internal/models"
	"github.com/jak103/powerplay/internal/server/apis"
	"github.com/jak103/powerplay/internal/server/services/auth"
)

// Pseudo database of messages
var messages []models.ChatMessage

func init() {
	apis.RegisterHandler(fiber.MethodGet, "/chat", auth.Public, getMessages)
	apis.RegisterHandler(fiber.MethodPost, "/chat", auth.Public, sendMessage)
}

func sendMessage(context *fiber.Ctx) error {
	var message models.ChatMessage
	err := context.BodyParser(&message)

	if err != nil {
		return err
	}

	messages = append(messages, message)

	return context.JSON(fiber.Map{
		"status": "success",
		"sent":   "true",
		"data":   message,
	})
}

func getMessages(context *fiber.Ctx) error {
	userID := context.Params("userID")
	var missedMessages []models.ChatMessage

	for _, message := range messages {
		if message.RecipientID == userID {
			missedMessages = append(missedMessages, message)
		}
	}

	return context.JSON(missedMessages)
}
