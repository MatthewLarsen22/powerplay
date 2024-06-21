package chat

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jak103/powerplay/internal/db"
	"github.com/jak103/powerplay/internal/models"
	"github.com/jak103/powerplay/internal/server/apis"
	"github.com/jak103/powerplay/internal/server/services/auth"
	"github.com/jak103/powerplay/internal/utils/log"
	"github.com/jak103/powerplay/internal/utils/responder"
)

var nextTagID = 0
var tags = make(map[string]TagConfiguration)

func init() {
	apis.RegisterHandler(fiber.MethodGet, "/chat/conversations", auth.Public, listConversations)
	apis.RegisterHandler(fiber.MethodPost, "/chat/conversations/create", auth.Public, createConversation)
	apis.RegisterHandler(fiber.MethodDelete, "/chat/conversations/delete", auth.Public, deleteConversation)
	apis.RegisterHandler(fiber.MethodPut, "/chat/conversations/updateName", auth.Public, updateName)
	apis.RegisterHandler(fiber.MethodPut, "/chat/conversations/updateimage", auth.Public, updateImage)
	apis.RegisterHandler(fiber.MethodPut, "/chat/conversations/updatedescription", auth.Public, updateDescription)
	apis.RegisterHandler(fiber.MethodPut, "/chat/conversations/adduser", auth.Public, addUser)
	apis.RegisterHandler(fiber.MethodPut, "/chat/conversations/removeuser", auth.Public, removeUser)
	apis.RegisterHandler(fiber.MethodPost, "/chat/tags/create", auth.Public, createTag)
}

func listConversations(c *fiber.Ctx) error {
	session := db.GetSession(c)
	conversationList, err := session.GetConversations()

	if err != nil {
		log.WithErr(err).Alert("Failed to parse conversation request payload")
		return responder.InternalServerError(c)
	}

	return responder.OkWithData(c, conversationList)
}

func createConversation(c *fiber.Ctx) error {
	conversationConfig := new(ConversationConfiguration)

	// Load the request body as a ConversationConfiguration object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(conversationConfig); err != nil {
		return responder.BadRequest(c, "Failed to parse the request. Please ensure that all inputs are of the correct data types, as specified in the API documentation.")
	}

	// Verify that values were provided for required fields. If any required values are missing, the request is bad.
	var errorMsg string
	if conversationConfig.Name == "" {
		if errorMsg != "" {
			errorMsg += "; "
		}
		errorMsg += "missing required field 'name'"
	}
	if conversationConfig.Type == "" {
		if errorMsg != "" {
			errorMsg += "; "
		}
		errorMsg += "missing required field 'type'"
	} else {
		if errorMsg != "" {
			errorMsg += "; "
		}
		if conversationConfig.Type != models.CHANNEL && conversationConfig.Type != models.DM {
			errorMsg += "invalid input for 'type'"
		}
	}
	if conversationConfig.Participants == nil {
		if errorMsg != "" {
			errorMsg += "; "
		}
		errorMsg += "missing required field 'participants'"
	}
	if errorMsg != "" {
		return responder.BadRequest(c, fmt.Sprintf("Failed to create the conversation. Please address the following issues with the request: %v.", errorMsg))
	}

	// Create a conversation using the provided data
	conversation := new(models.Conversation)
	conversation.Name = conversationConfig.Name
	conversation.Type = conversationConfig.Type
	conversation.Description = conversationConfig.Description
	conversation.ImageString = conversationConfig.ImageString
	conversation.Participants = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(conversationConfig.Participants)), ", "), "[]")

	// Save the new conversation to the database
	session := db.GetSession(c)
	record, err := session.CreateConversation(conversation)

	// If there was a problem saving the validated conversation to the database, report a server error
	if err != nil || record == nil {
		return responder.InternalServerError(c, fmt.Sprintf("There was a problem saving the new conversation to the database: %v.", err))
	}

	return responder.Ok(c)
}

func deleteConversation(c *fiber.Ctx) error {
	conversationID := new(ConversationID)

	// Load the request body as a conversationID string. If the provided value is the wrong type, the request is bad.
	if err := c.BodyParser(conversationID); err != nil {
		return responder.BadRequest(c, "Failed to parse the request. Please ensure that all inputs are of the correct data types, as specified in the API documentation.")
	}

	// Verify that a conversation ID was provided. If there is no conversation ID, the request is bad.
	var errorMsg string
	if conversationID.Value == 0 {
		if errorMsg != "" {
			errorMsg += "; "
		}
		errorMsg += "missing required field 'conversation_id'"
	}
	if errorMsg != "" {
		return responder.BadRequest(c, fmt.Sprintf("Failed to delete the conversation. Please address the following issues with the request: %v.", errorMsg))
	}

	// Delete the conversation from the database.
	session := db.GetSession(c)
	err := session.DeleteConversation(conversationID.Value)

	// If there was a problem deleting the conversation from the database, report a server error
	if err != nil {
		return responder.InternalServerError(c, fmt.Sprintf("There was a problem removing the conversation from the database: %v.", err))
	}

	return responder.Ok(c)
}

func updateImage(c *fiber.Ctx) error {
	updateData := new(ConversationPropertyChange)

	// Load the request body as a ConversationUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c, "Failed to parse the request. Please ensure that all inputs are of the correct data types, as specified in the API documentation.")
	}

	// Verify that values were provided for required fields and verify the existence of the conversation. If any required values are missing or the conversation doesn't exist, the request is bad.
	var errorMsg string
	if updateData.ConversationID == 0 {
		if errorMsg != "" {
			errorMsg += "; "
		}
		errorMsg += "missing required field 'conversation_id'"
	}
	if updateData.Value == "" {
		errorMsg += "missing required field 'value' is a required field.\n"
	}
	if errorMsg != "" {
		return responder.BadRequest(c, fmt.Sprintf("Failed to update conversation image. Please address the following issues with the request: %v.", errorMsg))
	}

	// Update the record in the database
	session := db.GetSession(c)
	err := session.UpdateConversationImage(updateData.ConversationID, updateData.Value)

	// If there was a problem updating the database record, report a server error
	if err != nil {
		return responder.InternalServerError(c, fmt.Sprintf("There was a problem updating the conversation image in the database: %v.", err))
	}

	return responder.Ok(c)
}

func updateDescription(c *fiber.Ctx) error {
	updateData := new(ConversationPropertyChange)

	// Load the request body as a ConversationUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c, "Failed to parse the request. Please ensure that all inputs are of the correct data types, as specified in the API documentation.")
	}

	// Verify that values were provided for required fields and verify the existence of the conversation. If any required values are missing or the conversation doesn't exist, the request is bad.
	var errorMsg string
	if updateData.ConversationID == 0 {
		if errorMsg != "" {
			errorMsg += "; "
		}
		errorMsg += "missing required field 'conversation_id'"
	}
	if updateData.Value == "" {
		if errorMsg != "" {
			errorMsg += "; "
		}
		errorMsg += "missing required field'value'"
	}
	if errorMsg != "" {
		return responder.BadRequest(c, fmt.Sprintf("Failed to update the conversation description. Please address the following issues with the request: %v.", errorMsg))
	}

	// Update the record in the database
	session := db.GetSession(c)
	err := session.UpdateConversationDescription(updateData.ConversationID, updateData.Value)

	// If there was a problem updating the database record, report a server error
	if err != nil {
		return responder.InternalServerError(c, fmt.Sprintf("There was a problem updating the description in the database: %v.", err))
	}

	return responder.Ok(c)
}

func updateName(c *fiber.Ctx) error {
	updateData := new(ConversationPropertyChange)

	// Load the request body as a ConversationUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c, "Failed to parse the request. Please ensure that all inputs are of the correct data types, as specified in the API documentation.")
	}

	// Verify that values were provided for required fields and verify the existence of the conversation. If any required values are missing or the conversation doesn't exist, the request is bad.
	var errorMsg string
	if updateData.ConversationID == 0 {
		if errorMsg != "" {
			errorMsg += "; "
		}
		errorMsg += "missing required field 'conversation_id'"
	}
	if updateData.Value == "" {
		if errorMsg != "" {
			errorMsg += "; "
		}
		errorMsg += "missing required field 'value'"
	}
	if errorMsg != "" {
		return responder.BadRequest(c, fmt.Sprintf("Failed to update conversation name. Please address the following issues with the request: %v.", errorMsg))
	}

	// Update the record in the database
	session := db.GetSession(c)
	err := session.UpdateConversationName(updateData.ConversationID, updateData.Value)

	// If there was a problem updating the database record, report a server error
	if err != nil {
		return responder.InternalServerError(c, fmt.Sprintf("There was a problem updating the name in the database: %v.", err))
	}

	return responder.Ok(c)
}

func addUser(c *fiber.Ctx) error {
	updateData := new(ConversationUserChange)

	// Load the request body as a ConversationUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c, "Failed to parse the request. Please ensure that all inputs are of the correct data types, as specified in the API documentation.")
	}

	// Verify that values were provided for required fields and verify the existence of the conversation. If any required values are missing or the conversation doesn't exist, the request is bad.
	var errorMsg string
	if updateData.ConversationID == 0 {
		if errorMsg != "" {
			errorMsg += "; "
		}
		errorMsg += "missing required field 'conversation_id'"
	}
	if updateData.UserID == 0 {
		if errorMsg != "" {
			errorMsg += "; "
		}
		errorMsg += "missing required field 'user_id'"
	}
	if errorMsg != "" {
		return responder.BadRequest(c, fmt.Sprintf("Failed to add the user to the conversation. Please address the following issues with the request: %v.", errorMsg))
	}

	// Update the Conversation record in the database with a new participants list
	session := db.GetSession(c)
	err := session.AddConversationParticipant(updateData.ConversationID, updateData.UserID)

	// If there was a problem updating the participants list in the database, report a server error
	if err != nil {
		return responder.InternalServerError(c, fmt.Sprintf("There was a problem updating the participants list in the database: %v.", err))
	}

	return responder.Ok(c)
}

func removeUser(c *fiber.Ctx) error {
	updateData := new(ConversationUserChange)

	// Load the request body as a ConversationUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c, "Failed to parse the request. Please ensure that all inputs are of the correct data types, as specified in the API documentation.")
	}

	// Verify that values were provided for required fields and verify the existence of the conversation. If any required values are missing or the conversation doesn't exist, the request is bad.
	var errorMsg string
	// var userIndex = -1
	if updateData.ConversationID == 0 {
		if errorMsg != "" {
			errorMsg += "; "
		}
		errorMsg += "missing required field 'conversation_id'"
	}
	if updateData.UserID == 0 {
		errorMsg += "missing required field 'value'"
	}
	if errorMsg != "" {
		return responder.BadRequest(c, fmt.Sprintf("Failed to remove the user from the conversation. Please address the following issues with the request: %v.", errorMsg))
	}

	// Update the Conversation record in the database with a new participants list
	session := db.GetSession(c)
	err := session.RemoveConversationParticipant(updateData.ConversationID, updateData.UserID)

	// If there was a problem updating the participants list in the database, report a server error
	if err != nil {
		return responder.InternalServerError(c, fmt.Sprintf("There was a problem updating the participants list in the database: %v.", err))
	}

	return responder.Ok(c)
}

func createTag(c *fiber.Ctx) error {
	tag := new(TagConfiguration)

	// Load the request body as a ChannelConfiguration object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(tag); err != nil {
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields. If any required values are missing, the request is bad.
	var errorMsg string
	if tag.Name == "" {
		errorMsg += "\t'name' is a required field.\n"
	}
	if tag.Description == "" {
		errorMsg += "\t'description' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("Tag creation failed. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	// Create a tag using the provided data
	tags[strconv.Itoa(nextTagID)] = *tag // TODO: store tags in the DB instead of just in a dictionary.
	nextTagID += 1
	log.Info("Tag created: " + tag.Name)
	return responder.Ok(c)
}

type ConversationID struct {
	Value uint `json:"conversation_id"`
}

type ConversationUserChange struct {
	ConversationID uint `json:"conversation_id"`
	UserID         uint `json:"user_id"`
}

type ConversationPropertyChange struct {
	ConversationID uint   `json:"conversation_id"`
	Value          string `json:"value"`
}

type ConversationConfiguration struct {
	Name         string                  `json:"name"`
	Type         models.ConversationType `json:"type"`
	Description  string                  `json:"description"`
	ImageString  string                  `json:"image_string"`
	Participants []uint                  `json:"participants"`
}

type TagID struct {
	Value string `json:"tag_id"`
}

type TagConfiguration struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
