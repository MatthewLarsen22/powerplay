package models

type Channel struct {
	DbModel
	ImageString  string `json:"image_string"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Participants []User `json:"participants" gorm:"many2many:channels_users"`
}
