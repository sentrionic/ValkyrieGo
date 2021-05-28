package handler

import (
	"github.com/go-redis/redis/v8"
	"github.com/sentrionic/valkyrie/model"
)

type WsServer struct {
	clients        map[*Client]bool
	register       chan *Client
	unregister     chan *Client
	broadcast      chan []byte
	rooms          map[*Room]bool
	channelService model.ChannelService
	guildService   model.GuildService
	userService    model.UserService
	rds            *redis.Client
}

// NewWebsocketServer creates a new WsServer type
func NewWebsocketServer(c *Config) *WsServer {
	return &WsServer{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan []byte),
		rooms:          make(map[*Room]bool),
		channelService: c.ChannelService,
		guildService:   c.GuildService,
		userService:    c.UserService,
		rds:            c.Redis,
	}
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run() {
	for {
		select {

		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)

		case message := <-server.broadcast:
			server.broadcastToClients(message)
		}
	}
}

func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
	}
}

func (server *WsServer) broadcastToClients(message []byte) {
	for client := range server.clients {
		client.send <- message
	}
}

func (server *WsServer) broadcastToRoom(message []byte, room string) {
	if room := server.findRoomById(room); room != nil {
		room.publishRoomMessage(message)
	}
}

func (server *WsServer) findRoomById(id string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetId() == id {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) createRoom(id string) *Room {
	room := NewRoom(id, server.rds)
	go room.RunRoom()
	server.rooms[room] = true

	return room
}
