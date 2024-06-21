package db

import (
	"github.com/jak103/powerplay/internal/models"
)

func (s session) CreateTag(tag *models.Tag) (*models.Tag, error) {
	result := s.connection.Create(tag)
	return resultOrError(tag, result)
}
