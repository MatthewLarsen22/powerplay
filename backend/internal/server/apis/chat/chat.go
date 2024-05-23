package chat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jak103/powerplay/internal/server/apis"
	"github.com/jak103/powerplay/internal/server/services/auth"
	"github.com/jak103/powerplay/internal/utils/log"
	"github.com/jak103/powerplay/internal/utils/responder"
)

var channels []ChannelConfiguration

func init() {
	apis.RegisterHandler(fiber.MethodGet, "/hello", auth.Public, helloWorld)
	apis.RegisterHandler(fiber.MethodPost, "/chat/channels/create", auth.Public, createChannel)
}

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}

func createChannel(c *fiber.Ctx) error {
	channel := new(ChannelConfiguration)

	// Load the request body as a ChannelConfiguration object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(channel); err != nil {
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields. If any required values are missing, the request is bad.
	var errorMsg string
	if channel.Name == "" {
		errorMsg = errorMsg + "\t'name' is a required field.\n"
	}
	if channel.MemberIDs == nil {
		errorMsg = errorMsg + "\t'member_ids' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("Channel creation failed. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	// Create a channel using the provided data
	channels = append(channels, *channel) // TODO: store channels in the DB instead of just in a slice.
	log.Info("Channel created: " + channel.Name)
	return responder.Ok(c)
}

type ChannelConfiguration struct {
	Name        string   `json:"name"`
	MemberIDs   []string `json:"member_ids"`
	ImageString string   `json:"image_string"`
	Description string   `json:"description"`
}
