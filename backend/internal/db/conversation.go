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
