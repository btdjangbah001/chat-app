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
	GroupID    uint    `json:"group_id"`
	Status     Receipt `json:"status"`
}

type UnsentAcknowledgement struct {
	MessageID       uint    `json:"message_id"`
	ReceiverID      uint    `json:"receiver_id"`
	GroupID         uint    `json:"group_id"`
	Status          Receipt `json:"status"`
	MessageSenderID uint    `json:"message_sender_id"`
}

func (unsentAcknowledgement *UnsentAcknowledgement) CreateUnsentAcknowledgement() (err error) {
	err = DB.Create(&unsentAcknowledgement).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUnsentAcknowledgementsForUser(userID uint) (*[]UnsentAcknowledgement, error) {
	var acknowledgements []UnsentAcknowledgement
	err := DB.Where("message_sender_id = ?", userID).Find(&acknowledgements).Error
	if err != nil {
		return nil, err
	}
	return &acknowledgements, nil
}

func DeleteUnsentAcknowledgementsForUser(userID uint) error {
	err := DB.Where("message_sender_id = ?", userID).Delete(&UnsentAcknowledgement{}).Error
	if err != nil {
		return err
	}
	return nil
}
