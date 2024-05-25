package models

import "time"

type ChatMessage struct {
	DbModel
	MessageID		uint		'json:"id" db: "id"' // could be a string for UUID
	SenderID		uint		'json:"senderID" db:"senderID"'
	Content			string		'json:"content" db:"content"'
	CreatedAt		time.Time	'json:"createdAt" db:"createdAt"'
	LastModifiedAt	time.Time	'json:"lastModifiedAt" db:"lastModifiedAt'
	DeletedAt		time.Time	'json:"deletedAt" db:"deletedAt"'
	RecipientID		uint		'json:"recipientID" db:"recipientID"'
	AttachmentURLs	[]string	'json:"attachmentURLs" db:"attachmentURLs"'
	Reactions		[]string	'json:"reactions" db:"reactions"' // assuming emojis are strings
	ParentID		uint		'json:"parentID" db:"parentID"'
}