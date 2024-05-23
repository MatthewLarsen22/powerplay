package chat

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jak103/powerplay/internal/server/apis"
	"github.com/jak103/powerplay/internal/server/services/auth"
	"github.com/jak103/powerplay/internal/utils/log"
	"github.com/jak103/powerplay/internal/utils/responder"
)

var nextID int = 0
var channels = make(map[string]ChannelConfiguration)

func init() {
	apis.RegisterHandler(fiber.MethodGet, "/hello", auth.Public, helloWorld)
	apis.RegisterHandler(fiber.MethodPost, "/chat/channels/create", auth.Public, createChannel)
	apis.RegisterHandler(fiber.MethodPut, "/chat/channels/adduser", auth.Public, addUser)
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
		errorMsg += "\t'name' is a required field.\n"
	}
	if channel.MemberIDs == nil {
		errorMsg += "\t'member_ids' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("Channel creation failed. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	// Create a channel using the provided data
	channels[strconv.Itoa(nextID)] = *channel // TODO: store channels in the DB instead of just in a dictionary.
	nextID += 1
	log.Info("Channel created: " + channel.Name)
	return responder.Ok(c)
}

func addUser(c *fiber.Ctx) error {
	updateData := new(ChannelUpdate)

	// Load the request body as a ChannelUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields and verify the existence of the channel. If any required values are missing or the channel doesn't exist, the request is bad.
	var errorMsg string
	var channel ChannelConfiguration
	var channelOk bool
	if updateData.ChannelID == "" {
		errorMsg += "\t'channel_id' is a required field.\n"
	} else {
		// Retrieve the channel specified in the channel_id field.
		channel, channelOk = channels[updateData.ChannelID] // TODO: retrieve the channel from the DB.
		if !channelOk {
			errorMsg += "\tNo channel exists with ID " + updateData.ChannelID + ".\n"
		}
	}
	if updateData.Value == "" {
		errorMsg += "\t'value' is a required field.\n"
	}
	// TODO: verify that the provided value is a valid user ID
	// TODO: verify that the user ID has not already been added to the channel
	if errorMsg != "" {
		log.Info("The user could not be added to the channel. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	channel.MemberIDs = append(channel.MemberIDs, updateData.Value) // TODO: update the member_ids list in the database
	log.Info("User " + updateData.Value + " added to channel " + updateData.ChannelID)
	return responder.Ok(c)
}

type ChannelUpdate struct {
	ChannelID string `json:"channel_id"`
	Value     string `json:"value"`
}

type ChannelConfiguration struct {
	Name        string   `json:"name"`
	MemberIDs   []string `json:"member_ids"`
	ImageString string   `json:"image_string"`
	Description string   `json:"description"`
}
