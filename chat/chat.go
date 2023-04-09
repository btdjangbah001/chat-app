package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/btdjangbah001/chat-app/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var connections = make(map[uint]*websocket.Conn)

var connectionsMutex sync.Mutex

func addConnection(userId uint, conn *websocket.Conn) {
	// Lock the mutex to synchronize access to the connections map
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	// Add the connection to the map of active connections
	connections[userId] = conn

	// Set a close handler for the WebSocket connection to remove it from the map when it is closed
	conn.SetCloseHandler(func(code int, text string) error {
		// Lock the mutex to synchronize access to the connections map
		connectionsMutex.Lock()
		defer connectionsMutex.Unlock()

		// Remove the connection from the map of active connections
		delete(connections, userId)

		return nil
	})
}

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
	// defer ws.Close()

	val, exists := c.Get("user")
	if !exists {
		_ = fmt.Errorf("error getting user from context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := val.(*models.User)

	// Add the connection to the map of connections
	addConnection(user.ID, ws)
	// user.Ws = ws

	if err = sendUnreadMessages(user); err != nil {
		_ = fmt.Errorf("error sending unread messages: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "something went wrong please try again"})
		return
	}

	// Create a channel for incoming messages
	messageChan := make(chan []byte)

	// Start a goroutine to read incoming messages from the WebSocket connection and send them to the channel
	go func() {
		for {
			_, messageBytes, err := ws.ReadMessage()
			if err != nil {
				_ = fmt.Errorf("error reading message from WebSocket: %v", err)
				// Handle the error
				continue
			}
			fmt.Printf(string(messageBytes))

			messageChan <- messageBytes
		}
	}()

	// Start a goroutine to handle incoming messages from the channel
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
			fmt.Printf(message.Content)

			switch message.Type {
			case models.PRIVATE:
				fmt.Printf("private")
				err := sendMessage(message.RecipientID, &message)
				if err != nil {
					_ = fmt.Errorf("error sending message: %v", err)
					// Handle the error
					continue
				}

			case models.GROUP:
				// Send the message to all connections that belong to the group
				groupParticipants, err := models.GetGroupParticipants(message.RecipientID)
				if err != nil {
					_ = fmt.Errorf("error getting group participants: %v", err)
					// Handle the error
					continue
				}

				for _, participant := range *groupParticipants {
					err := sendMessage(participant, &message)
					if err != nil {
						// Handle the error
						continue
					}
				}
			}
		}
	}()
}

func sendUnreadMessages(user *models.User) error {
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

		if err := sendMessage(user.ID, &message); err != nil {
			return err
		}
	}
	return nil
}

func keepUnsentMessages(message *models.UnsentMessage) error {
	err := message.CreateUnsentMessage()
	if err != nil {
		return err
	}
	return nil
}

func sendMessage(recipient_id uint, message *models.Message) error {
	// recipient, err := models.GetUser(message.RecipientID)
	// if err != nil {
	// 	// Handle the error
	// 	return err
	// }
	var ws *websocket.Conn

	ws, ok := connections[recipient_id]
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
