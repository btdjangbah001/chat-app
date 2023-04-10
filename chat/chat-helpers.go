package chat

import (
	"github.com/btdjangbah001/chat-app/models"
	"github.com/gorilla/websocket"
)

func keepUnsentMessages(message *models.UnsentMessage) error {
	err := message.CreateUnsentMessage()
	if err != nil {
		return err
	}
	return nil
}

func SendMessage(recipient_id uint, message *models.Message) error {
	var ws *websocket.Conn

	ws, ok := Connections[recipient_id]
	if !ok {
		unsentMessage := models.UnsentMessage{
			Content:        message.Content,
			RecipientID:    message.RecipientID,
			SenderID:       message.SenderID,
			SenderUsername: message.SenderUsername,
			Type:           message.Type,
		}
		keepUnsentMessages(&unsentMessage)
	} else {
		clientMessage := models.ClientMessage{
			Content:        message.Content,
			SenderUsername: message.SenderUsername,
			Type:           message.Type,
			GroupID:        message.GroupID,
		}
		err := ws.WriteJSON(clientMessage)
		if err != nil {
			// Handle the error
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
