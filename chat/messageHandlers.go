package chat

import (
	"encoding/json"
	"fmt"

	"github.com/btdjangbah001/chat-app/models"
)

func HandleMessage(msg *IncomingMessage) error {
	var message models.Message
	err := json.Unmarshal(msg.Data, &message)
	if err != nil {
		_ = fmt.Errorf("error unmarshalling message: %v", err)
		// Handle the error
		return err
	}

	// Send acknowlegdement that the message was received
	ack := models.Acknowledgement{
		MessageID:  message.MessageID,
		ReceiverID: message.RecipientID,
		GroupID:    message.GroupID,
		Status:     models.SENT,
	}

	if err = SendAcknowledgement(&ack, message.SenderID); err != nil {
		_ = fmt.Errorf("error sending acknowlegdement: %v", err)
		// Handle the error
		return err
	}

	switch message.Type {
	case models.PRIVATE:
		err := SendMessage(message.RecipientID, &message)
		if err != nil {
			_ = fmt.Errorf("error sending message: %v", err)
			// Handle the error
			return err
		}

	case models.GROUP:
		// Send the message to all connections that belong to the group
		groupParticipants, err := models.GetGroupParticipantsExceptUser(message.GroupID, message.SenderID)
		if err != nil {
			_ = fmt.Errorf("error getting group participants: %v", err)
			// Handle the error
			return err
		}

		for _, participant := range *groupParticipants {
			message.RecipientID = participant
			err := SendMessage(participant, &message)
			if err != nil {
				// Handle the error
				break
			}
		}
	}

	return nil
}

func HandleAcknowlegdement(msg *IncomingMessage) error {
	var ack models.Acknowledgement
	err := json.Unmarshal(msg.Data, &ack)
	if err != nil {
		_ = fmt.Errorf("error unmarshalling acknowlegement: %v", err)
		// Handle the error
		return err
	}

	SendAcknowledgement(&ack, ack.ReceiverID)

	return nil
}
