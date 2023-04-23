package chat

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

var Connections = make(map[uint]*websocket.Conn)

var connectionsMutex sync.Mutex

func AddConnection(userId uint, conn *websocket.Conn) {
	// Lock the mutex to synchronize access to the connections map
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	// Add the connection to the map of active connections
	Connections[userId] = conn

	// Set a close handler for the WebSocket connection to remove it from the map when it is closed
	conn.SetCloseHandler(func(code int, text string) error {
		// Lock the mutex to synchronize access to the connections map
		connectionsMutex.Lock()
		defer connectionsMutex.Unlock()

		fmt.Printf("Connection closed for user %d\n", userId)
		// Remove the connection from the map of active connections
		delete(Connections, userId)

		return nil
	})
}
