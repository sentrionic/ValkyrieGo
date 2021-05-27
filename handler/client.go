package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/handler/ws"
	"github.com/sentrionic/valkyrie/model"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var newline = []byte{'\n'}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == os.Getenv("CORS_ORIGIN")
	},
}

// Client represents the websocket client at the server
type Client struct {
	// The actual websocket connection.
	ID       string
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	rooms    map[*Room]bool
}

func newClient(conn *websocket.Conn, wsServer *WsServer, id string) *Client {
	return &Client{
		ID:       id,
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    make(map[*Room]bool),
	}
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)

	_ = client.conn.SetReadDeadline(time.Now().Add(pongWait))

	client.conn.SetPongHandler(func(string) error {
		_ = client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		client.handleNewMessage(jsonMessage)
	}

}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				_ = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	client.wsServer.unregister <- client
	for room := range client.rooms {
		room.unregister <- client
	}
	close(client.send)
	_ = client.conn.Close()
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(wsServer *WsServer, ctx *gin.Context) {

	userId := ctx.MustGet("userId").(string)
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, wsServer, userId)

	go client.writePump()
	go client.readPump()

	wsServer.register <- client
}

func (client *Client) handleNewMessage(jsonMessage []byte) {

	var message model.ReceivedMessage
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
	}

	switch message.Action {
	// Join Room Actions
	case ws.JoinChannelAction:
		client.handleJoinChannelMessage(message)
	case ws.JoinGuildAction:
		client.handleJoinGuildMessage(message)
	case ws.JoinUserAction:
		client.handleJoinRoomMessage(message)

	// Leave Room Actions
	case ws.LeaveRoomAction:
		client.handleLeaveRoomMessage(message)
	case ws.LeaveGuildAction:
		client.handleLeaveGuildMessage(message)

	// Chat Typing Actions
	case ws.StartTypingAction:
		roomID := message.Room
		if room := client.wsServer.findRoomById(roomID); room != nil {
			msg := model.WebsocketMessage{
				Action: ws.AddToTypingAction,
				Data:   message.Message,
			}
			room.broadcast <- &msg
		}

	case ws.StopTypingAction:
		roomID := message.Room
		if room := client.wsServer.findRoomById(roomID); room != nil {
			msg := model.WebsocketMessage{
				Action: ws.RemoveFromTypingAction,
				Data:   message.Message,
			}
			room.broadcast <- &msg
		}

	// Online Status Actions
	case ws.ToggleOnlineAction:
	case ws.ToggleOfflineAction:

	// Other
	case ws.GetRequestCountAction:
	}
}

func (client *Client) handleJoinChannelMessage(message model.ReceivedMessage) {
	roomName := message.Room

	cs := client.wsServer.channelService
	channel, err := cs.Get(roomName)

	if err != nil {
		return
	}

	if err = cs.IsChannelMember(channel, client.ID); err != nil {
		return
	}

	client.handleJoinRoomMessage(message)
}

func (client *Client) handleJoinGuildMessage(message model.ReceivedMessage) {
	roomName := message.Room

	gs := client.wsServer.guildService
	guild, err := gs.GetGuild(roomName)

	if err != nil {
		return
	}

	if !isMember(guild, client.ID) {
		return
	}

	client.handleJoinRoomMessage(message)
}

func (client *Client) handleJoinRoomMessage(message model.ReceivedMessage) {
	roomName := message.Room

	room := client.wsServer.findRoomById(roomName)
	if room == nil {
		room = client.wsServer.createRoom(roomName)
	}

	client.rooms[room] = true

	room.register <- client
}

func (client *Client) handleLeaveGuildMessage(message model.ReceivedMessage) {
	_ = client.wsServer.guildService.UpdateMemberLastSeen(client.ID, message.Room)
	client.handleLeaveRoomMessage(message)
}

func (client *Client) handleLeaveRoomMessage(message model.ReceivedMessage) {
	room := client.wsServer.findRoomById(message.Room)
	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.unregister <- client
}
