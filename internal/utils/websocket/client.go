package websocket

import (
	"chs_cloud_general/internal/global_var"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var WSConn *websocket.Conn

type TServerStatus struct {
	DayendClose bool      `json:"dayend_close"`
	AuditDate   time.Time `json:"audit_date"`
	ServerTime  time.Time `json:"server_time"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
	c := s.conn
	defer func() {
		h.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		m := message{msg, s.room}
		h.broadcast <- m
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// write writes a message with the given message type and payload.
func (c *connection) writeJSON(payload interface{}) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteJSON(payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs() gin.HandlerFunc {
	return func(c *gin.Context) {
		global_var.MxSocket.Lock()
		go h.run()
		var upgrader = websocket.Upgrader{
			Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")}, // <-- add this line
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		clientID := c.Param("clientId")
		if err != nil {
			log.Println(err)
			return
		}
		// send server status on connect
		// IsRunning, _ := global_query.IsDayendCloseRunning(c)
		//send dayend status to client
		// DayendStatus, _ := global_query.GetDayendCloseStatus(c)

		con := &connection{send: make(chan []byte, 256), ws: ws}
		s := subscription{con, clientID}
		h.register <- s

		// SendInitialMessage(s, global_var.WSMessageType.Connection, nil, global_var.WSDataType.ServerStatus, TServerStatus{
		// 	DayendClose: IsRunning,
		// 	ServerTime:  time.Now(),
		// 	AuditDate:   global_query.GetAuditDate(c,DB),
		// })
		// SendMessage(global_var.UserInfo.CompanyCode, global_var.WSMessageType.Connection, nil, global_var.WSDataType.DayendCloseStatus, DayendStatus)
		global_var.MxSocket.Unlock()

		go s.writePump()
		s.readPump()
	}
}
