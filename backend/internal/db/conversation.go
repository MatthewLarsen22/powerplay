package db

import (
	"fmt"
	"strings"

	"github.com/jak103/powerplay/internal/models"
)

func (s session) GetConversations() ([]models.Conversation, error) {
	conversations := make([]models.Conversation, 0)
	err := s.connection.Find(&conversations)
	return resultsOrError(conversations, err)
}

func (s session) GetConversation(conversationID uint) (*models.Conversation, error) {
	var conversation *models.Conversation
	err := s.connection.First(conversation, conversationID)
	return resultOrError(conversation, err)
}

func (s session) CreateConversation(conversation *models.Conversation) (*models.Conversation, error) {
	result := s.connection.Create(conversation)
	return resultOrError(conversation, result)
}

func (s session) DeleteConversation(conversationID uint) error {
	resultObj := s.connection.Delete(&models.Conversation{}, conversationID)
	if resultObj.Error == nil && resultObj.RowsAffected < 1 {
		resultObj.Error = fmt.Errorf("could not find a record for conversation %v", conversationID)
	}
	return resultObj.Error
}

func (s session) UpdateConversationName(conversationID uint, newName string) error {
	resultObj := s.connection.Model(&models.Conversation{}).Where("id = ?", conversationID).Update("name", newName)
	if resultObj.Error == nil && resultObj.RowsAffected < 1 {
		resultObj.Error = fmt.Errorf("could not find a record for conversation %v", conversationID)
	}
	return resultObj.Error
}

func (s session) UpdateConversationDescription(conversationID uint, newDescription string) error {
	resultObj := s.connection.Model(&models.Conversation{}).Where("id = ?", conversationID).Update("description", newDescription)
	if resultObj.Error == nil && resultObj.RowsAffected < 1 {
		resultObj.Error = fmt.Errorf("could not find a record for conversation %v", conversationID)
	}
	return resultObj.Error
}

func (s session) UpdateConversationImage(conversationID uint, newImageString string) error {
	resultObj := s.connection.Model(&models.Conversation{}).Where("id = ?", conversationID).Update("image_string", newImageString)
	if resultObj.Error == nil && resultObj.RowsAffected < 1 {
		resultObj.Error = fmt.Errorf("could not find a record for conversation %v", conversationID)
	}
	return resultObj.Error
}

func (s session) RemoveConversationParticipant(conversationID uint, userID uint) error {
	var conversation *models.Conversation
	resultObj := s.connection.Select("participants").First(&conversation).Where("id = ?", conversationID)

	if resultObj.Error == nil && conversation == nil {
		resultObj.Error = fmt.Errorf("could not find a record for conversation %v", conversationID)
	}

	if resultObj.Error == nil {
		splitUsers := strings.Split(conversation.Participants, ", ")

		participantIndex := -1
		for index, participantID := range splitUsers {
			if participantID == fmt.Sprintf("%v", userID) {
				participantIndex = index
				break
			}
		}

		if participantIndex < 0 {
			resultObj.Error = fmt.Errorf("user %v is not a participant in conversation %v", userID, conversationID)
		} else {
			splitUsers = append(splitUsers[:participantIndex], splitUsers[participantIndex+1:]...)
			conversation.Participants = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(splitUsers)), ", "), "[]")
			resultObj = s.connection.Model(&models.Conversation{}).Where("id = ?", conversationID).Update("participants", conversation.Participants)
		}
	}

	return resultObj.Error
}

func (s session) AddConversationParticipant(conversationID uint, userID uint) error {
	var conversation *models.Conversation
	resultObj := s.connection.Select("participants").First(&conversation).Where("id = ?", conversationID)

	if resultObj.Error == nil && conversation == nil {
		resultObj.Error = fmt.Errorf("could not find a record for conversation %v", conversationID)
	}

	if resultObj.Error == nil {
		var user *models.User
		userResult := s.connection.First(user, userID)

		if userResult.Error != nil {
			resultObj.Error = userResult.Error
		} else {
			if user == nil {
				resultObj.Error = fmt.Errorf("could not find a record for user %v", userID)
			}
		}
	}

	if resultObj.Error == nil {
		splitUsers := strings.Split(conversation.Participants, ", ")

		participantIndex := -1
		for index, participantID := range splitUsers {
			if participantID == fmt.Sprintf("%v", userID) {
				participantIndex = index
				break
			}
		}

		if participantIndex < 0 {
			splitUsers = append(splitUsers, fmt.Sprintf("%v", userID))
			conversation.Participants = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(splitUsers)), ", "), "[]")
			resultObj = s.connection.Model(&models.Conversation{}).Where("id = ?", conversationID).Update("participants", conversation.Participants)
		} else {
			resultObj.Error = fmt.Errorf("user %v is already a participant in conversation %v", userID, conversationID)
		}
	}

	return resultObj.Error
}
