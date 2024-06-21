package models

type Tag struct {
	DbModel
	TagName     string `json:"tag_name"`
	Description string `json:"description"`
}

type UserTag struct {
	DbModel
	TagID  uint `json:"tag_id"`
	UserID uint `json:"user_id"`
}
