package models

type Receipt int

const (
	SENT Receipt = iota
	RECEIVED
	READ
	ERROR
)

type Acknowledgement struct {
	MessageID  uint    `json:"message_id"`
	ReceiverID uint    `json:"receiver_id"`
	Status     Receipt `json:"status"`
}

type UnsentAcknowledgement struct {
	MessageID  uint    `json:"message_id"`
	ReceiverID uint    `json:"receiver_id"`
	Status     Receipt `json:"status"`
}
