package models

type Channel struct {
	DbModel
	ImageURL    string   `json:"image_url"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MemberIDs   []string `json:"member_ids" gorm:"type:text[]"`
}
