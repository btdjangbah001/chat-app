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

	if err = SendAcknowledgement(&ack, ack.ReceiverID); err != nil {
		_ = fmt.Errorf("error sending acknowlegdement: %v", err)
		// Handle the error
		return err
	}

	return nil
}

func HandleStatus(msg *IncomingMessage) error {
	var status models.Status
	err := json.Unmarshal(msg.Data, &status)
	if err != nil {
		_ = fmt.Errorf("error unmarshalling status: %v", err)
		// Handle the error
		return err
	}

	outMsg := OutgoingMessage{
		Type: "status",
	}

	switch status.Activity {
	case models.OFFLINE:
	case models.ONLINE:
		if err = HandleUserOnlineOrOfflineStatus(&status, &outMsg); err != nil {
			return err
		}
	case models.THINKING:
	case models.TYPING:
		if err = HandleTypingThinkingStatus(&status, &outMsg); err != nil {
			return err
		}
	}

	return nil
}

func HandleUserOnlineOrOfflineStatus(status *models.Status, out *OutgoingMessage) error {
	data := models.Status{
		SenderID:   status.SenderID,
		ReceiverID: status.ReceiverID,
	}
	_, ok := Connections[status.ReceiverID]
	if !ok {
		data.Activity = models.OFFLINE
	} else {
		data.Activity = models.ONLINE
	}
	out.Data = data
	if err := SendStatus(out, status.SenderID); err != nil {
		_ = fmt.Errorf("error sending status: %v", err)
		return err
	}
	return nil
}

func HandleTypingThinkingStatus(status *models.Status, out *OutgoingMessage) error {
	ws, ok := Connections[status.ReceiverID]
	if !ok {
		return nil
	}
	data := models.Status{
		SenderID:   status.SenderID,
		ReceiverID: status.ReceiverID,
	}
	if status.Activity == models.TYPING {
		data.Activity = models.TYPING
	} else {
		data.Activity = models.THINKING
	}
	out.Data = data
	if err := ws.WriteJSON(out); err != nil {
		_ = fmt.Errorf("error sending typing status: %v", err)
		return err
	}
	return nil
}
