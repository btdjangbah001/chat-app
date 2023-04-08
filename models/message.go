package models

type Type int

const (
	PRIVATE Type = iota
	GROUP
)

type Message struct {
	Content     string `json:"content"`
	SenderID    uint `json:"sender_id"`
	RecipientID uint `json:"recipient_id"`
	Type        Type   `json:"type"`
	CreatedAt   string `json:"created_at"`
}

type UnsentMessage struct {
	Content     string `json:"content"`
	SenderID    uint `json:"sender_id"`
	RecipientID uint `json:"recipient_id"`
	Type        Type   `json:"type"`
	GroupID     string `json:"group_id"`
}

type ClientMessage struct {
	Content     string `json:"content"`
	SenderID    uint `json:"sender_id"`
	Type		Type   `json:"type"`
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
	err := DB.Where("recipient_id = ? AND read = ?", userID, false).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return &messages, nil
}