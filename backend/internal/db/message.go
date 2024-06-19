package db

import "github.com/jak103/powerplay/internal/models"

func (s session) GetMessages() ([]models.Message, error) {
	messages := make([]models.Message, 0)
	err := s.connection.Find(&messages)
	return resultsOrError(messages, err)
}

func (s session) CreateMessage(message *models.Message) (*models.Message, error) {
	result := s.connection.Create(message)
	return resultOrError(message, result)
}

func (s session) DeleteMessage(messageID uint) error {
	result := s.connection.Delete(&models.Message{}, messageID)
	return result.Error
}

func (s session) UpdateMessageContent(messageID uint, newContent string) error {
	resultObj := s.connection.Model(&models.Message{}).Where("id = ?", messageID).Update("content", newContent)
	return resultObj.Error
}
