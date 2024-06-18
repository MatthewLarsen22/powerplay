package models

type ConversationType string

const (
	TEAM    ConversationType = "team"
	LEAGUE  ConversationType = "league"
	DM      ConversationType = "dm"
	CHANNEL ConversationType = "channel"
)

type Conversation struct {
	DbModel
	Name         string           `json:"name"`
	Type         ConversationType `json:"type"`
	Description  string           `json:"description"`
	ImageString  string           `json:"image_string"`
	Participants []uint           `json:"participants" gorm:"type:bigint[]"`
}
