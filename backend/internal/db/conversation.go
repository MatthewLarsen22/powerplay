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
	result := s.connection.Delete(&models.Conversation{}, conversationID)
	return result.Error
}

func (s session) UpdateConversationName(conversationID uint, newName string) error {
	resultObj := s.connection.Model(&models.Conversation{}).Where("id = ?", conversationID).Update("name", newName)
	return resultObj.Error
}

func (s session) UpdateConversationDescription(conversationID uint, newDescription string) error {
	resultObj := s.connection.Model(&models.Conversation{}).Where("id = ?", conversationID).Update("description", newDescription)
	return resultObj.Error
}

func (s session) UpdateConversationImage(conversationID uint, newImageString string) error {
	resultObj := s.connection.Model(&models.Conversation{}).Where("id = ?", conversationID).Update("image_string", newImageString)
	return resultObj.Error
}

func (s session) RemoveConversationParticipant(conversationID uint, userID uint) error {
	var conversation *models.Conversation
	resultObj := s.connection.Select("participants").First(&conversation).Where("id = ?", conversationID)
	if resultObj.Error == nil {
		splitUsers := strings.Split(conversation.Participants, ", ")

		participantIndex := -1
		for index, participantID := range splitUsers {
			if participantID == fmt.Sprintf("%v", userID) {
				participantIndex = index
				break
			}
		}

		if participantIndex >= 0 {
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
		}
	}

	return resultObj.Error
}
