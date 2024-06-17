package chat

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jak103/powerplay/internal/db"
	"github.com/jak103/powerplay/internal/models"
	"github.com/jak103/powerplay/internal/server/apis"
	"github.com/jak103/powerplay/internal/server/services/auth"
	"github.com/jak103/powerplay/internal/utils/log"
	"github.com/jak103/powerplay/internal/utils/responder"
)

func init() {
	apis.RegisterHandler(fiber.MethodGet, "/hello", auth.Public, helloWorld)
	apis.RegisterHandler(fiber.MethodGet, "/chat/conversations", auth.Public, listConversations)
	apis.RegisterHandler(fiber.MethodPost, "/chat/conversations/create", auth.Public, createConversation)
	apis.RegisterHandler(fiber.MethodDelete, "/chat/conversations/delete", auth.Public, deleteConversation)
	apis.RegisterHandler(fiber.MethodPut, "/chat/conversations/updateimage", auth.Public, updateImage)
	apis.RegisterHandler(fiber.MethodPut, "/chat/conversations/updatedescription", auth.Public, updateDescription)
	apis.RegisterHandler(fiber.MethodPut, "/chat/conversations/rename", auth.Public, rename)
	apis.RegisterHandler(fiber.MethodPut, "/chat/conversations/adduser", auth.Public, addUser)
	apis.RegisterHandler(fiber.MethodPut, "/chat/conversations/removeuser", auth.Public, removeUser)
}

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello World")
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
	log.Info("Creating a conversation")
	conversationConfig := new(ConversationConfiguration)

	// Load the request body as a ConversationConfiguration object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(conversationConfig); err != nil {
		log.Info("Encountered an error while parsing the body of the request")
		log.Info(err.Error())
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields. If any required values are missing, the request is bad.
	var errorMsg string
	if conversationConfig.Name == "" {
		errorMsg += "\t'name' is a required field.\n"
	}
	if conversationConfig.Type == "" {
		errorMsg += "\t'type' is  required field.\n"
	} else {
		if conversationConfig.Type != models.CHANNEL && conversationConfig.Type != models.DM {
			errorMsg += "\t'type' must be \"channel\" or \"dm\"."
		}
	}
	if conversationConfig.Participants == nil {
		errorMsg += "\t'participants' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("Conversation creation failed. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
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
		log.WithErr(err).Alert("There was a problem saving the new conversation to the database")
		return responder.InternalServerError(c)
	}

	log.Info("Conversation created: " + conversation.Name)

	return responder.Ok(c)
}

func deleteConversation(c *fiber.Ctx) error {
	conversationID := new(ConversationID)

	// Load the request body as a conversationID string. If the provided value is the wrong type, the request is bad.
	if err := c.BodyParser(conversationID); err != nil {
		return responder.BadRequest(c)
	}

	// Verify that a conversation ID was provided. If there is no conversation ID, the request is bad.
	var errorMsg string
	if conversationID.Value == 0 {
		errorMsg += "\t'conversation_id' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("The conversation could not be deleted. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	// Delete the conversation from the database.
	session := db.GetSession(c)
	err := session.DeleteConversation(conversationID.Value)

	// If there was a problem deleting the conversation from the database, report a server error
	if err != nil {
		log.Info("There was a problem removing the conversation from the database.")
		return responder.InternalServerError(c)
	}

	log.Info("Deleted conversation")

	return responder.Ok(c)
}

func updateImage(c *fiber.Ctx) error {
	updateData := new(ConversationPropertyChange)

	// Load the request body as a ConversationUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields and verify the existence of the conversation. If any required values are missing or the conversation doesn't exist, the request is bad.
	var errorMsg string
	if updateData.ConversationID == 0 {
		errorMsg += "\t'conversation_id' is a required field.\n"
	}
	if updateData.Value == "" {
		errorMsg += "\t'value' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("Conversation image could not be updated. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	// Update the record in the database
	session := db.GetSession(c)
	err := session.UpdateConversationImage(updateData.ConversationID, updateData.Value)

	// If there was a problem updating the database record, report a server error
	if err != nil {
		log.Info("There was a problem updating the conversation image in the database.")
		return responder.InternalServerError(c)
	}

	log.Info("Conversation image updated")

	return responder.Ok(c)
}

func updateDescription(c *fiber.Ctx) error {
	updateData := new(ConversationPropertyChange)

	// Load the request body as a ConversationUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields and verify the existence of the conversation. If any required values are missing or the conversation doesn't exist, the request is bad.
	var errorMsg string
	if updateData.ConversationID == 0 {
		errorMsg += "\t'conversation_id' is a required field.\n"
	}
	if updateData.Value == "" {
		errorMsg += "\t'value' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("Conversation description could not be updated. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	// Update the record in the database
	session := db.GetSession(c)
	err := session.UpdateConversationDescription(updateData.ConversationID, updateData.Value)

	// If there was a problem updating the database record, report a server error
	if err != nil {
		log.Info("There was a problem updating the description in the database.")
		return responder.InternalServerError(c)
	}

	log.Info("Conversation description updated")

	return responder.Ok(c)
}

func rename(c *fiber.Ctx) error {
	updateData := new(ConversationPropertyChange)

	// Load the request body as a ConversationUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields and verify the existence of the conversation. If any required values are missing or the conversation doesn't exist, the request is bad.
	var errorMsg string
	if updateData.ConversationID == 0 {
		errorMsg += "\t'conversation_id' is a required field.\n"
	}
	if updateData.Value == "" {
		errorMsg += "\t'value' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("Conversation description could not be updated. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	// Update the record in the database
	session := db.GetSession(c)
	err := session.UpdateConversationName(updateData.ConversationID, updateData.Value)

	// If there was a problem updating the database record, report a server error
	if err != nil {
		log.Info("There was a problem updating the name in the database.")
		return responder.InternalServerError(c)
	}

	log.Info("Conversation name updated")

	return responder.Ok(c)
}

func addUser(c *fiber.Ctx) error {
	updateData := new(ConversationUserChange)

	// Load the request body as a ConversationUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields and verify the existence of the conversation. If any required values are missing or the conversation doesn't exist, the request is bad.
	var errorMsg string
	if updateData.ConversationID == 0 {
		errorMsg += "\t'conversation_id' is a required field.\n"
	}
	if updateData.UserID == 0 {
		errorMsg += "\t'user_id' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("The user could not be added to the conversation. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	// Update the Conversation record in the database with a new participants list
	session := db.GetSession(c)
	err := session.AddConversationParticipant(updateData.ConversationID, updateData.UserID)

	// If there was a problem updating the participants list in the database, report a server error
	if err != nil {
		log.Info("There was a problem updating the participants list in the database.")
		return responder.InternalServerError(c)
	}

	log.Info("User added to conversation")

	return responder.Ok(c)
}

func removeUser(c *fiber.Ctx) error {
	updateData := new(ConversationUserChange)

	// Load the request body as a ConversationUpdate object. If any of the provided values are the wrong type, the request is bad.
	if err := c.BodyParser(updateData); err != nil {
		return responder.BadRequest(c)
	}

	// Verify that values were provided for required fields and verify the existence of the conversation. If any required values are missing or the conversation doesn't exist, the request is bad.
	var errorMsg string
	// var userIndex = -1
	if updateData.ConversationID == 0 {
		errorMsg += "\t'conversation_id' is a required field.\n"
	}
	if updateData.UserID == 0 {
		errorMsg += "\t'value' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("The user could not be removed from the conversation. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	session := db.GetSession(c)
	err := session.RemoveConversationParticipant(updateData.ConversationID, updateData.UserID)

	if err != nil {
		log.Info("There was a problem updating the participants list in the database.")
		return responder.InternalServerError(c)
	}

	log.Info("User removed from conversation.")

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
