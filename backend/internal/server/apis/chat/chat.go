package chat

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/jak103/powerplay/internal/models"
	"github.com/jak103/powerplay/internal/server/apis"

	"github.com/jak103/powerplay/internal/utils/log"
)

var (
	messages  []models.ChatMessage // psuedo-database
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan models.ChatMessage)
)

func init() {
	apis.RegisterHandler(fiber.MethodGet, "/chat", nil, websocket.New(HandleConnections))

	go HandleMessages()
}

func WebsocketHandler(context *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(context) {
		context.Locals("allowed", true)
		return context.Next()
	}

	return fiber.ErrUpgradeRequired
}

func HandleConnections(connection *websocket.Conn) {
	defer func() {
		delete(clients, connection)
		connection.Close()
	}()

	clients[connection] = true

	for {
		var message models.ChatMessage
		err := connection.ReadJSON(&message)

		if err != nil {
			log.WithErr(err).Alert("Cannot Read message...")
			delete(clients, connection)
			break
		}

		message.CreatedAt = time.Now()
		broadcast <- message
	}
}

func HandleMessages() {
	for {
		message := <-broadcast
		messages = append(messages, message)

		for client := range clients {
			err := client.WriteJSON(message)

			if err != nil {
				log.WithErr(err).Alert("Cannot write message")
				client.Close()
				delete(clients, client)
			}
		}
	}
}
