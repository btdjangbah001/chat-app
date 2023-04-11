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

	if err = SendUnreadMessages(user); err != nil {
		_ = fmt.Errorf("error sending unread messages: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "something went wrong please try again"})
		return
	}

	// Create a channel for incoming messages
	messageChan := make(chan []byte)

	// Read incoming messages from the WebSocket connection and send them to the channel
	go func() {
		for {
			_, messageBytes, err := ws.ReadMessage()
			if err != nil {
				_ = fmt.Errorf("error reading message from WebSocket: %v", err)
				// Handle the error
				break
			}

			messageChan <- messageBytes
		}
	}()

	// Handle incoming messages from the channel
	go func() {
		for messageBytes := range messageChan {
			// Process the message (send to the appropriate recipients, store in the database, etc.)
			var message models.Message
			err = json.Unmarshal(messageBytes, &message)
			if err != nil {
				_ = fmt.Errorf("error unmarshalling message: %v", err)
				// Handle the error
				continue
			}

			switch message.Type {
			case models.PRIVATE:
				fmt.Printf("private")
				err := SendMessage(message.RecipientID, &message)
				if err != nil {
					_ = fmt.Errorf("error sending message: %v", err)
					// Handle the error
					continue
				}

			case models.GROUP:
				// Send the message to all connections that belong to the group
				groupParticipants, err := models.GetGroupParticipantsExceptUser(message.GroupID, message.SenderID)
				if err != nil {
					_ = fmt.Errorf("error getting group participants: %v", err)
					// Handle the error
					continue
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
		}
	}()
}
