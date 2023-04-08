package models

type Type int

const (
	PRIVATE Type = iota
	GROUP
)

type Message struct {
	Content     string `json:"content"`
	SendorID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
	Type        Type   `json:"type"`
	CreatedAt   string `json:"created_at"`
}
