package db

import "github.com/jak103/powerplay/internal/models"

func (s session) GetConversations() ([]models.Conversation, error) {
	conversations := make([]models.Conversation, 0)
	err := s.connection.Find(&conversations)
	return resultsOrError(conversations, err)
}

func (s session) CreateConversation(conversation *models.Conversation) (*models.Conversation, error) {
	result := s.connection.Create(conversation)
	return resultOrError(conversation, result)
}

func (s session) UpdateConversationName(conversationID uint, newName string) error {
	resultObj := s.connection.Model(&models.Conversation{}).Where("id = ?", conversationID).Update("description", newName)
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
