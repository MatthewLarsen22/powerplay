package models

import (
	"time"
)

type Message struct {
	ID        		uint        `json:"id" gorm:"primarykey"`
	CreatedAt 		time.Time   `json:"created_at"`
	SenderID  		uint 		`json:"senderID"`
	Content   		string 	 	`json:"content"`
	LastModifiedAt  time.Time 	`json:"lastModifiedAt"`
	DeletedAt 		*time.Time  `json:"deleted_at,omitempty" gorm:"index"`
	RecipientID 	uint 		`json:"recipientID"`
	AttachmentURLS  []string 	`json:"attachmentURLS"`
	Reactions 		[]string 	`json:"reactions"` 						
	// I'm not sure what reactions would be stored as, maybe ints if they're like ASCII characters. shrug
}