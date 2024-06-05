package models

type ChannelType string

const (
	TEAM    ChannelType = "team"
	LEAGUE  ChannelType = "league"
	DM      ChannelType = "dm"
	CHANNEL ChannelType = "channel"
)

type Channel struct {
	DbModel
	Name         string      `json:"name"`
	Type         ChannelType `json:"type"`
	Description  string      `json:"description"`
	ImageString  string      `json:"image_string"`
	Participants []uint      `json:"participants" gorm:"type:bigint[]"`
}
