package models

type StatusType int

const (
	ONLINE StatusType = iota
	TYPING
	THINKING
	OFFLINE
)

type Status struct {
	Activity   StatusType `json:"activity_id"`
	ReceiverID uint       `json:"receiver_id"`
	GroupID    uint       `json:"group_id"`
	SenderID   uint       `json:"sender_id"`
}
