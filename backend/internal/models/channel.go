package models

type Channel struct {
	DbModel
	ImageString string `json:"image_string"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MemberIDs   []uint `json:"member_ids" gorm:"type:integer[]"`
}
