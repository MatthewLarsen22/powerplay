package models

type Channel struct {
	DbModel
	ImageString string   `json:"image_string"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MemberIDs   []string `json:"member_ids" gorm:"type:text[]"`
}
