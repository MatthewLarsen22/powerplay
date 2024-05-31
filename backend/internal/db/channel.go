package db

import "github.com/jak103/powerplay/internal/models"

func (s session) GetChannels() ([]models.Channel, error) {
	channels := make([]models.Channel, 0)
	err := s.connection.Find(&channels)
	return resultsOrError(channels, err)
}

func (s session) CreateChannel(request *models.Channel) error {
	result := s.connection.Create(request)
	return result.Error
}
