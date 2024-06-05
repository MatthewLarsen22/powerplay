package models

type Channel struct {
	DbModel
	ImageString  string `json:"image_string"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Participants []uint `json:"participants" gorm:"type:bigint[]"`
}
