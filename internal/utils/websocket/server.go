package websocket

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsCon struct {
	conn     *websocket.Conn
	clientId string
}

var (
	upgrader = websocket.Upgrader{
		Subprotocols: []string{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODQyNDEyNjEsInJlZnJlc2giOmZhbHNlLCJ1c2VyIjoiU1lTVEVNIn0.p5c8Fj7oT5XxWq70d5ryZFXorV2bZedpJDIDlKYTInQ"}, // <-- add this line
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	clients = make(map[string]map[*websocket.Conn]bool)
)

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	clientId := c.Param("clientId")
	if err != nil {
		fmt.Println("Failed to upgrade WebSocket connection:", err)
		return
	}

	connections := clients[clientId]
	if connections == nil {
		connections = make(map[*websocket.Conn]bool)
		clients[clientId] = connections
	}
	clients[clientId][conn] = true

	//send server status on connectt
	// IsRunning, _ := global_query.IsDayendCloseRunning(c)
	// SendInitialMessage(conn, global_var.WSMessageType.Connection, nil, global_var.WSDataType.ServerStatus, TServerStatus{
	// 	DayendClose: IsRunning,
	// 	ServerTime:  time.Now(),
	// 	AuditDate:   global_query.GetAuditDate(c,DB),
	// })

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Failed to read message from WebSocket connection:", err)
			delete(clients, clientId)
			break
		}

		fmt.Println("Received message from WebSocket connection:", string(message))

		for client := range clients[clientId] {
			fmt.Println(client)
			if client != conn {
				err := client.WriteMessage(messageType, message)
				if err != nil {
					fmt.Println("Failed to write message to WebSocket connection:", err)
					delete(clients, clientId)
				}
			}
		}
	}
}

func SendMessage(ClientId string, messageType interface{}, message interface{}, dataType, data interface{}, userId string) error {
	connections := h.rooms[ClientId]
	for c := range connections {
		fmt.Println(c)
		err := c.writeJSON(gin.H{
			"user_id":      userId,
			"message_type": messageType,
			"message":      message,
			"data_type":    dataType,
			"data":         data,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func SendInitialMessage(s subscription, messageType interface{}, message interface{}, dataType, data interface{}) error {
	err := s.conn.writeJSON(gin.H{
		"message_type": messageType,
		"message":      message,
		"data_type":    dataType,
		"data":         data,
	})
	if err != nil {
		return err
	}
	return nil
}
