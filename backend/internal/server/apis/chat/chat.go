package chat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jak103/powerplay/internal/db"
	"github.com/jak103/powerplay/internal/models"
	"github.com/jak103/powerplay/internal/server/apis"
	"github.com/jak103/powerplay/internal/server/services/auth"
	"github.com/jak103/powerplay/internal/utils/log"
	"github.com/jak103/powerplay/internal/utils/responder"
	"github.com/lib/pq"
)

var nextID uint = 1
var channels = make(map[uint]models.Channel)

func init() {
	apis.RegisterHandler(fiber.MethodGet, "/hello", auth.Public, helloWorld)
	apis.RegisterHandler(fiber.MethodGet, "/chat/channels", auth.Public, listChannels)
	apis.RegisterHandler(fiber.MethodPost, "/chat/channels/create", auth.Public, createChannel)
	apis.RegisterHandler(fiber.MethodDelete, "/chat/channels/delete", auth.Public, deleteChannel)
	apis.RegisterHandler(fiber.MethodPut, "/chat/channels/updateimage", auth.Public, updateImage)
	apis.RegisterHandler(fiber.MethodPut, "/chat/channels/updatedescription", auth.Public, updateDescription)
	apis.RegisterHandler(fiber.MethodPut, "/chat/channels/rename", auth.Public, rename)
	apis.RegisterHandler(fiber.MethodPut, "/chat/channels/adduser", auth.Public, addUser)
	apis.RegisterHandler(fiber.MethodPut, "/chat/channels/removeuser", auth.Public, removeUser)
}

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}

func listChannels(c *fiber.Ctx) error {
	session := db.GetSession(c)
	channelList, err := session.GetChannels()

	if err != nil {
		log.WithErr(err).Alert("Failed to parse channel request payload")
		return responder.InternalServerError(c)
	}

	return responder.OkWithData(c, channelList)
}

func createChannel(c *fiber.Ctx) error {
	log.Info("Creating a channel")
	channelConfig := new(ChannelConfiguration)

	// Load the request body as a ChannelConfiguration object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(channelConfig); err != nil {
		log.Info("Encountered an error while parsing the body of the request")
		log.Info(err.Error())
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields. If any required values are missing, the request is bad.
	var errorMsg string
	if channelConfig.Name == "" {
		errorMsg += "\t'name' is a required field.\n"
	}
	if channelConfig.Participants == nil {
		errorMsg += "\t'participants' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("Channel creation failed. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	// Create a channel using the provided data

	channel := new(models.Channel)
	channel.Name = channelConfig.Name
	channel.Description = channelConfig.Description
	channel.ImageString = channelConfig.ImageString
	channel.Participants = pq.Int64(channelConfig.Participants)
	session := db.GetSession(c)
	record, err := session.CreateChannel(channel)

	channels[nextID] = *channel // TODO: store channels in the DB instead of just in a dictionary.
	log.Info("Channel created: " + channel.Name)
	nextID += 1

	if err != nil {
		log.WithErr(err).Alert("Failed to parse channel request payload")
		log.Info("Participants: %T", channel.Participants)
		return responder.InternalServerError(c)
	}

	if record == nil {
		return responder.BadRequest(c, "Could not post channel into database")
	}

	return responder.Ok(c)
}

func deleteChannel(c *fiber.Ctx) error {
	channelID := new(ChannelID)

	// Load the request body as a channelID string. If the provided value is the wrong type, the request is bad.
	if err := c.BodyParser(channelID); err != nil {
		return responder.BadRequest(c)
	}

	// Verify the existence of the channel. If the channel doesn't exist, the request is bad.
	var errorMsg string
	_, channelExists := channels[channelID.Value] // TODO: retrieve the channel from the DB
	if !channelExists {
		errorMsg += "\tChannel does not exist.\n"
	}
	if errorMsg != "" {
		log.Info("The channel could not be deleted. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	delete(channels, channelID.Value) // TODO: delete the channel from the DB
	log.Info("Deleted channel")
	return responder.Ok(c)
}

func updateImage(c *fiber.Ctx) error {
	updateData := new(ChannelPropertyChange)

	// Load the request body as a ChannelUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields and verify the existence of the channel. If any required values are missing or the channel doesn't exist, the request is bad.
	var errorMsg string
	var channel models.Channel
	var channelOk bool
	if updateData.ChannelID == 0 {
		errorMsg += "\t'channel_id' is a required field.\n"
	} else {
		// Retrieve the channel specified in the channel_id field.
		channel, channelOk = channels[updateData.ChannelID] // TODO: retrieve the channel from the DB.
		if !channelOk {
			errorMsg += "\tSpecified channel does not exist.\n"
		}
	}
	if updateData.Value == "" {
		errorMsg += "\t'value' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("Channel image could not be updated. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	channel.ImageString = updateData.Value // TODO: update the participants list in the database
	channels[updateData.ChannelID] = channel
	log.Info("Channel image updated")
	return responder.Ok(c)
}

func updateDescription(c *fiber.Ctx) error {
	updateData := new(ChannelPropertyChange)

	// Load the request body as a ChannelUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields and verify the existence of the channel. If any required values are missing or the channel doesn't exist, the request is bad.
	var errorMsg string
	var channel models.Channel
	var channelOk bool
	if updateData.ChannelID == 0 {
		errorMsg += "\t'channel_id' is a required field.\n"
	} else {
		// Retrieve the channel specified in the channel_id field.
		channel, channelOk = channels[updateData.ChannelID] // TODO: retrieve the channel from the DB.
		if !channelOk {
			errorMsg += "\tThe specified channel does not exist.\n"
		}
	}
	if updateData.Value == "" {
		errorMsg += "\t'value' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("Channel description could not be updated. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	channel.Description = updateData.Value // TODO: update the participants list in the database
	channels[updateData.ChannelID] = channel
	log.Info("Channel description updated")
	return responder.Ok(c)
}

func rename(c *fiber.Ctx) error {
	updateData := new(ChannelPropertyChange)

	// Load the request body as a ChannelUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields and verify the existence of the channel. If any required values are missing or the channel doesn't exist, the request is bad.
	var errorMsg string
	var channel models.Channel
	var channelOk bool
	if updateData.ChannelID == 0 {
		errorMsg += "\t'channel_id' is a required field.\n"
	} else {
		// Retrieve the channel specified in the channel_id field.
		channel, channelOk = channels[updateData.ChannelID] // TODO: retrieve the channel from the DB.
		if !channelOk {
			errorMsg += "\tThe specified channel does not exist.\n"
		}
	}
	if updateData.Value == "" {
		errorMsg += "\t'value' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("Channel description could not be updated. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	channel.Name = updateData.Value
	channels[updateData.ChannelID] = channel
	log.Info("Channel name updated")
	return responder.Ok(c)
}

func addUser(c *fiber.Ctx) error {
	updateData := new(ChannelUserChange)

	// Load the request body as a ChannelUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields and verify the existence of the channel. If any required values are missing or the channel doesn't exist, the request is bad.
	var errorMsg string
	var channel models.Channel
	var channelOk bool
	if updateData.ChannelID == 0 {
		errorMsg += "\t'channel_id' is a required field.\n"
	} else {
		// Retrieve the channel specified in the channel_id field.
		channel, channelOk = channels[updateData.ChannelID] // TODO: retrieve the channel from the DB.
		if !channelOk {
			errorMsg += "\tThe specified channel does not exist.\n"
		}
	}
	if updateData.UserID == 0 {
		errorMsg += "\t'user_id' is a required field.\n"
	}
	// TODO: verify that the provided value is a valid user ID
	// TODO: verify that the user ID has not already been added to the channel
	if errorMsg != "" {
		log.Info("The user could not be added to the channel. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	// TODO: retrieve the user corresponding to the ID user from the database
	user := models.User{
		FirstName: "anonymous",
	}
	user.ID = updateData.UserID

	// channel.Participants = append(channel.Participants, user) // TODO: update the participants list in the database
	channels[updateData.ChannelID] = channel
	log.Info("User added to channel")
	return responder.Ok(c)
}

func removeUser(c *fiber.Ctx) error {
	updateData := new(ChannelUserChange)

	// Load the request body as a ChannelUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields and verify the existence of the channel. If any required values are missing or the channel doesn't exist, the request is bad.
	var errorMsg string
	var channel models.Channel
	var channelOk bool
	var userIndex = -1
	if updateData.ChannelID == 0 {
		errorMsg += "\t'channel_id' is a required field.\n"
	} else {
		// Retrieve the channel specified in the channel_id field.
		channel, channelOk = channels[updateData.ChannelID] // TODO: retrieve the channel from the DB.
		if !channelOk {
			errorMsg += "\tThe specified channel does not exist.\n"
		}
	}
	if updateData.UserID == 0 {
		errorMsg += "\t'value' is a required field.\n"
	} else {
		if channelOk {
			// for i, v := range channel.Participants {
			// 	if v.ID == updateData.UserID {
			// 		userIndex = i
			// 		break
			// 	}
			// }

			if userIndex < 0 {
				errorMsg += "\tThe specified user is not a participant in the specified channel.\n"
			}
		}
	}
	if errorMsg != "" {
		log.Info("The user could not be removed from the channel. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	// TODO: update the participants list in the database
	// if len(channel.Participants) > 1 {
	// 	channel.Participants = append(channel.Participants[:userIndex], channel.Participants[userIndex+1:]...)
	// } else {
	// 	channel.Participants = make([]models.User, 0)
	// }
	channels[updateData.ChannelID] = channel
	log.Info("User removed from channel.")
	return responder.Ok(c)
}

type ChannelID struct {
	Value uint `json:"channel_id"`
}

type ChannelUserChange struct {
	ChannelID uint `json:"channel_id"`
	UserID    uint `json:"user_id"`
}

type ChannelPropertyChange struct {
	ChannelID uint   `json:"channel_id"`
	Value     string `json:"value"`
}

type ChannelConfiguration struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ImageString  string `json:"image_string"`
	Participants []uint `json:"participants"`
}
