package chat

import (
	"encoding/json"

	"github.com/btdjangbah001/chat-app/models"
)

type IncomingMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type OutgoingMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func keepUnsentMessages(message *models.UnsentMessage) error {
	err := message.CreateUnsentMessage()
	if err != nil {
		return err
	}
	return nil
}

func SendMessage(recipient_id uint, message *models.Message) error {
	ws, ok := Connections[recipient_id]
	if !ok {
		unsentMessage := models.UnsentMessage{
			Content:        message.Content,
			MessageID:      message.MessageID,
			RecipientID:    message.RecipientID,
			SenderID:       message.SenderID,
			SenderUsername: message.SenderUsername,
			Type:           message.Type,
			GroupID:        message.GroupID,
		}
		keepUnsentMessages(&unsentMessage)
	} else {
		msg := models.ClientMessage{
			Content:        message.Content,
			MessageID:      message.MessageID,
			SenderUsername: message.SenderUsername,
			Type:           message.Type,
			GroupID:        message.GroupID,
		}

		clientMessage := OutgoingMessage{
			Type: "message",
			Data: msg,
		}

		err := ws.WriteJSON(clientMessage)
		if err != nil {
			// Handle the error
			return err
		}
		if err = SendAcknowledgement(&models.Acknowledgement{
			MessageID:  message.MessageID,
			ReceiverID: message.RecipientID,
			Status:     models.RECEIVED,
		}, message.SenderID); err != nil {
			return err
		}
	}
	return nil
}

func SendUnreadMessages(user *models.User) error {
	// Check the message queue for unsent messages
	unsentMessages, _ := models.GetUnreadMessagesForUser(user.ID)

	// Send any unsent messages to the user
	for _, unsentMessage := range *unsentMessages {
		message := models.Message{
			Content:        unsentMessage.Content,
			RecipientID:    unsentMessage.RecipientID,
			SenderID:       unsentMessage.SenderID,
			Type:           unsentMessage.Type,
			SenderUsername: unsentMessage.SenderUsername,
		}

		if err := SendMessage(user.ID, &message); err != nil {
			return err
		}
	}

	err := models.DeleteUnreadMessagesForUser(user.ID)

	if err != nil {
		return err
	}

	return nil
}

func SendAcknowledgement(acknowledgement *models.Acknowledgement, receiver_id uint) error {
	ws, ok := Connections[receiver_id]
	if !ok {
		if errr := keepUnsentAcknowledgements(&models.UnsentAcknowledgement{
			ReceiverID:      acknowledgement.ReceiverID,
			MessageID:       acknowledgement.MessageID,
			GroupID:         acknowledgement.GroupID,
			Status:          acknowledgement.Status,
			MessageSenderID: receiver_id,
		}); errr != nil {
			return errr
		}
		return nil
	} else {
		err := ws.WriteJSON(acknowledgement)
		if err != nil {
			// Handle the error
			return err
		}
	}
	return nil
}

func keepUnsentAcknowledgements(acknowledgement *models.UnsentAcknowledgement) error {
	err := acknowledgement.CreateUnsentAcknowledgement()
	if err != nil {
		return err
	}
	return nil
}

func SendUnreadAcknowledgements(user *models.User) error {
	// Check the message queue for unsent messages
	unsentAcknowledgements, _ := models.GetUnsentAcknowledgementsForUser(user.ID)

	// Send any unsent messages to the user
	for _, unsentAcknowledgement := range *unsentAcknowledgements {
		acknowledgement := models.Acknowledgement{
			ReceiverID: unsentAcknowledgement.ReceiverID,
			MessageID:  unsentAcknowledgement.MessageID,
			GroupID:    unsentAcknowledgement.GroupID,
			Status:     unsentAcknowledgement.Status,
		}

		if err := SendAcknowledgement(&acknowledgement, user.ID); err != nil {
			return err
		}
	}

	err := models.DeleteUnsentAcknowledgementsForUser(user.ID)

	if err != nil {
		return err
	}

	return nil
}

func SendStatus(msg *OutgoingMessage, receiverID uint) error {
	ws, ok := Connections[receiverID]
	if !ok {
		return nil
	}
	err := ws.WriteJSON(msg)
	if err != nil {
		// Handle the error
		return err
	}
	return nil
}
