
package models

import "time"

type ChatMessage struct {
	DbModel
	MessageID 		string 		`json:"messageID" db:"messageID"`
	SenderID 		string 		`json:"senderID" db:"senderID"`
	Content 		string 		`json:"content" db:"content"`
	CreatedAt 		time.Time 	`json:"createdAt" db:"createdAt"`
	LastModifiedAt 	time.Time 	`json:"lastModifiedAt" db:"lastModifiedAt"`
	DeletedAt 		time.Time 	`json:"deletedAt" db:"deletedAt"`
	RecipientID 	string 		`json:"recipientID" db:"recipientID"`
	AttachmentUrls 	[]string 	`json:"attachmentUrls" db:"attachmentUrls"`
	Reactions	 	[]string 	`json:"reactions" db:"reactions"`
	ParentID	 	string 		`json:"parentID" db:"parentID"`
}