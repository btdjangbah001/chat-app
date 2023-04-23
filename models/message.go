package models

type Type int

const (
	PRIVATE Type = iota
	GROUP
)

type Message struct {
	Content        string `json:"content"`
	MessageID      uint   `json:"message_id"`
	SenderID       uint   `json:"sender_id"`
	SenderUsername string `json:"sender_username"`
	RecipientID    uint   `json:"recipient_id"`
	Type           Type   `json:"type"`
	CreatedAt      string `json:"created_at"`
	GroupID        uint   `json:"group_id"`
}

type UnsentMessage struct {
	Content        string `json:"content"`
	MessageID      uint   `json:"message_id"`
	SenderID       uint   `json:"sender_id"`
	SenderUsername string `json:"sender_username"`
	RecipientID    uint   `json:"recipient_id"`
	Type           Type   `json:"type"`
	GroupID        uint   `json:"group_id"`
}

type ClientMessage struct {
	Content        string `json:"content"`
	MessageID      uint   `json:"message_id"`
	SenderUsername string `json:"sender_username"`
	Type           Type   `json:"type"`
	GroupID        uint   `json:"group_id"`
}

func (unsentMessage *UnsentMessage) CreateUnsentMessage() (err error) {
	err = DB.Create(&unsentMessage).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUnreadMessagesForUser(userID uint) (*[]UnsentMessage, error) {
	var messages []UnsentMessage
	err := DB.Where("recipient_id = ?", userID).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return &messages, nil
}

func DeleteUnreadMessagesForUser(userID uint) error {
	err := DB.Where("recipient_id = ?", userID).Delete(&UnsentMessage{}).Error
	if err != nil {
		return err
	}
	return nil
}
