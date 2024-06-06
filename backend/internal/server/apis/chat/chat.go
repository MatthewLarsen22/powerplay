package chat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jak103/powerplay/internal/db"
	"github.com/jak103/powerplay/internal/models"
	"github.com/jak103/powerplay/internal/server/apis"
	"github.com/jak103/powerplay/internal/server/services/auth"
	"github.com/jak103/powerplay/internal/utils/log"
	"github.com/jak103/powerplay/internal/utils/responder"
)

var nextID uint = 1
var conversations = make(map[uint]models.Conversation)

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
	// conversation.Participants = conversationConfig.Participants
	session := db.GetSession(c)
	record, err := session.CreateConversation(conversation)

	conversations[nextID] = *conversation // TODO: store conversations in the DB instead of just in a dictionary.
	log.Info("Conversation created: " + conversation.Name)
	nextID += 1

	if err != nil {
		log.WithErr(err).Alert("Failed to parse conversation request payload")
		log.Info("Participants: %T", conversation.Participants)
		return responder.InternalServerError(c)
	}

	if record == nil {
		return responder.BadRequest(c, "Could not post conversation into database")
	}

	return responder.Ok(c)
}

func deleteConversation(c *fiber.Ctx) error {
	conversationID := new(ConversationID)

	// Load the request body as a conversationID string. If the provided value is the wrong type, the request is bad.
	if err := c.BodyParser(conversationID); err != nil {
		return responder.BadRequest(c)
	}

	// Verify the existence of the conversation. If the conversation doesn't exist, the request is bad.
	var errorMsg string
	_, conversationExists := conversations[conversationID.Value] // TODO: retrieve the conversation from the DB
	if !conversationExists {
		errorMsg += "\tConversation does not exist.\n"
	}
	if errorMsg != "" {
		log.Info("The conversation could not be deleted. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	delete(conversations, conversationID.Value) // TODO: delete the conversation from the DB
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
	var conversation models.Conversation
	var conversationOk bool
	if updateData.ConversationID == 0 {
		errorMsg += "\t'conversation_id' is a required field.\n"
	} else {
		// Retrieve the conversation specified in the conversation_id field.
		conversation, conversationOk = conversations[updateData.ConversationID] // TODO: retrieve the conversation from the DB.
		if !conversationOk {
			errorMsg += "\tSpecified conversation does not exist.\n"
		}
	}
	if updateData.Value == "" {
		errorMsg += "\t'value' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("Conversation image could not be updated. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	conversation.ImageString = updateData.Value // TODO: update the participants list in the database
	conversations[updateData.ConversationID] = conversation
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
	var conversation models.Conversation
	var conversationOk bool
	if updateData.ConversationID == 0 {
		errorMsg += "\t'conversation_id' is a required field.\n"
	} else {
		// Retrieve the conversation specified in the conversation_id field.
		conversation, conversationOk = conversations[updateData.ConversationID] // TODO: retrieve the conversation from the DB.
		if !conversationOk {
			errorMsg += "\tThe specified conversation does not exist.\n"
		}
	}
	if updateData.Value == "" {
		errorMsg += "\t'value' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("Conversation description could not be updated. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	conversation.Description = updateData.Value // TODO: update the participants list in the database
	conversations[updateData.ConversationID] = conversation
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
	var conversation models.Conversation
	var conversationOk bool
	if updateData.ConversationID == 0 {
		errorMsg += "\t'conversation_id' is a required field.\n"
	} else {
		// Retrieve the conversation specified in the conversation_id field.
		conversation, conversationOk = conversations[updateData.ConversationID] // TODO: retrieve the conversation from the DB.
		if !conversationOk {
			errorMsg += "\tThe specified conversation does not exist.\n"
		}
	}
	if updateData.Value == "" {
		errorMsg += "\t'value' is a required field.\n"
	}
	if errorMsg != "" {
		log.Info("Conversation description could not be updated. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	conversation.Name = updateData.Value
	conversations[updateData.ConversationID] = conversation
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
	var conversation models.Conversation
	var conversationOk bool
	if updateData.ConversationID == 0 {
		errorMsg += "\t'conversation_id' is a required field.\n"
	} else {
		// Retrieve the conversation specified in the conversation_id field.
		conversation, conversationOk = conversations[updateData.ConversationID] // TODO: retrieve the conversation from the DB.
		if !conversationOk {
			errorMsg += "\tThe specified conversation does not exist.\n"
		}
	}
	if updateData.UserID == 0 {
		errorMsg += "\t'user_id' is a required field.\n"
	}
	// TODO: verify that the provided value is a valid user ID
	// TODO: verify that the user ID has not already been added to the conversation
	if errorMsg != "" {
		log.Info("The user could not be added to the conversation. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	// TODO: retrieve the user corresponding to the ID user from the database
	user := models.User{
		FirstName: "anonymous",
	}
	user.ID = updateData.UserID

	// conversation.Participants = append(conversation.Participants, user) // TODO: update the participants list in the database
	conversations[updateData.ConversationID] = conversation
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
	var conversation models.Conversation
	var conversationOk bool
	var userIndex = -1
	if updateData.ConversationID == 0 {
		errorMsg += "\t'conversation_id' is a required field.\n"
	} else {
		// Retrieve the conversation specified in the conversation_id field.
		conversation, conversationOk = conversations[updateData.ConversationID] // TODO: retrieve the conversation from the DB.
		if !conversationOk {
			errorMsg += "\tThe specified conversation does not exist.\n"
		}
	}
	if updateData.UserID == 0 {
		errorMsg += "\t'value' is a required field.\n"
	} else {
		if conversationOk {
			// for i, v := range conversation.Participants {
			// 	if v.ID == updateData.UserID {
			// 		userIndex = i
			// 		break
			// 	}
			// }

			if userIndex < 0 {
				errorMsg += "\tThe specified user is not a participant in the specified conversation.\n"
			}
		}
	}
	if errorMsg != "" {
		log.Info("The user could not be removed from the conversation. Reason(s):\n" + errorMsg)
		return responder.BadRequest(c)
	}

	// TODO: update the participants list in the database
	// if len(conversation.Participants) > 1 {
	// 	conversation.Participants = append(conversation.Participants[:userIndex], conversation.Participants[userIndex+1:]...)
	// } else {
	// 	conversation.Participants = make([]models.User, 0)
	// }
	conversations[updateData.ConversationID] = conversation
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
