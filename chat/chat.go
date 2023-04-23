package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/btdjangbah001/chat-app/models"
	"github.com/btdjangbah001/chat-app/utilities"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func ChatHandler(c *gin.Context) {
	// Upgrade the HTTP request to a WebSocket connection
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Upgrade the TCP connection to a WebSocket connection
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		_ = fmt.Errorf("error upgrading request to WebSocket: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "something went wrong please try again"})
		return
	}

	user := utilities.GetLoggedInUser(c)

	// Add the connection to the map of connections
	AddConnection(user.ID, ws)

	if err = SendUnreadAcknowledgements(user); err != nil {
		_ = fmt.Errorf("error sending unread acknowledgements: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "something went wrong please try again"})
		return
	}

	if err = SendUnreadMessages(user); err != nil {
		_ = fmt.Errorf("error sending unread messages: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "something went wrong please try again"})
		return
	}

	// Create a channel for incoming messages
	messageChan := make(chan IncomingMessage)
	done := make(chan bool)

	// Read incoming messages from the WebSocket connection and send them to the channel
	go func() {
		for {
			fmt.Println("Message reading for user ", user.ID)
			var message IncomingMessage
			err := ws.ReadJSON(&message)
			if err != nil {
				_ = fmt.Errorf("error reading message from WebSocket: %v", err)
				fmt.Println("error reading message from WebSocket: ", err)
				close(done)
				break
			}

			messageChan <- message
		}
	}()

	// Handle incoming messages from the channel
	go func() {
		fmt.Println("Starting message handler for user ", user.ID)
		for {
			select {
			case <-done:
				fmt.Println("Closing message handler for user ", user.ID)
				close(messageChan)
				return
			case msg := <-messageChan:
				switch msg.Type {
				case "message":
					// Process the message (send to the appropriate recipients, store in the database, etc.)
					err := HandleMessage(&msg)
					if err != nil {
						continue
					}

				case "acknowlegdement":
					err = HandleAcknowlegdement(&msg)
					if err != nil {
						continue
					}
				}

			}
		}
	}()
}

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
