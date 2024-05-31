package models

import "time"

type ChatMessage struct {
	DbModel
	MessageID      string    `json:"messageID" db:"messageID"` // could be a string for UUID
	SenderID       string    `json:"senderID" db:"senderID"`
	Content        string    `json:"content" db:"content"`
	CreatedAt      time.Time `json:"createdAt" db:"createdAt"`
	LastModifiedAt time.Time `json:"lastModifiedAt" db:"lastModifiedAt"`
	DeletedAt      time.Time `json:"deletedAt" db:"deletedAt"`
	RecipientID    string    `json:"recipientID" db:"recipientID"`
	AttachmentURLs []string  `json:"attachmentURLs" db:"attachmentURLs"`
	Reactions      []string  `json:"reactions" db:"reactions"` // assuming emojis are strings
	ParentID       string    `json:"parentID" db:"parentID"`
}
