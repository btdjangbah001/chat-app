package chat

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/btdjangbah001/chat-app/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type connection struct {
	ws     *websocket.Conn // The WebSocket connection
	userId uint            // The ID of the user who made the connection
}

var connections = make(map[*connection]bool) // Map of all connections

var connectionsMutex sync.Mutex

func addConnection(conn *connection) {
	// Lock the mutex to synchronize access to the connections map
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	// Add the connection to the map of active connections
	connections[conn] = true

	// Set a close handler for the WebSocket connection to remove it from the map when it is closed
	conn.ws.SetCloseHandler(func(code int, text string) error {
		// Lock the mutex to synchronize access to the connections map
		connectionsMutex.Lock()
		defer connectionsMutex.Unlock()

		// Remove the connection from the map of active connections
		delete(connections, conn)

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
		// Handle error
		return
	}
	defer ws.Close()

	val, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := val.(*models.User)

	// Add the connection to the map of connections
	addConnection(&connection{ws, user.ID})
	user.Ws = ws

	err = models.UpdateUserWebsocket(user, ws)
	if err != nil {
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
				// Handle the error
				continue
			}

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
				// Handle the error
				continue
			}

			switch message.Type {
			case models.PRIVATE:
				err := sendMessage(message.RecipientID, &message)
				if err != nil {
					// Handle the error
					continue
				}

			case models.GROUP:
				// Send the message to all connections that belong to the group
				groupParticipants, err := models.GetGroupParticipants(message.RecipientID)
				if err != nil {
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

	// Logic for failed messages to be sent again will be implemented here

	// for {
	// 	//get messaged from the db that eblongs to the user and send it to the user

	// 	// Check the message queue for unsent messages
	// 	unsentMessages, _ := models.GetUnreadMessagesForUser(user.ID)

	// 	// If there are no unsent messages, wait for new messages to arrive on the channel
	// 	if len(*unsentMessages) == 0 {
	// 		select {
	// 		case message := <-messageChan:
	// 			// Process the message (send to the appropriate recipients, store in the database, etc.)
	// 		}
	// 	}

	// 	// Send any unsent messages to the user
	// 	for _, message := range unsentMessages {
	// 		err := conn.ws.WriteMessage(websocket.TextMessage, []byte(message))
	// 		if err != nil {
	// 			// Handle the error
	// 			break
	// 		}
	// 	}
	// }
}

func keepUnsentMessages(message *models.UnsentMessage) error {
	err := message.CreateUnsentMessage()
	if err != nil {
		return err
	}
	return nil
}

func sendMessage(recipient_id uint, message *models.Message) error {
	recipient, err := models.GetUser(message.RecipientID)
	if err != nil {
		// Handle the error
		return err
	}

	if _, ok := connections[&connection{recipient.Ws, recipient.ID}]; !ok {
		unsentMessage := models.UnsentMessage{
			Content:     message.Content,
			RecipientID: message.RecipientID,
			SenderID:    message.SenderID,
			Type:        message.Type,
		}
		keepUnsentMessages(&unsentMessage)
	} else {
		err = recipient.Ws.WriteJSON(message)
		if err != nil {
			// Handle the error
			return err
		}
	}
	return nil
}
