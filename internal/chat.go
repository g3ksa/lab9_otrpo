package internal

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/go-redis/redis"
)

var ChatChannel = os.Getenv("REDIS_CHAT_CHANNEL")

type Message struct {
	Type   string   `json:"type"`
	Author string   `json:"author,omitempty"`
	Text   string   `json:"text,omitempty"`
	Users  []string `json:"users,omitempty"`
}

type Hub struct {
	clients     map[*Client]bool
	register    chan *Client
	unregister  chan *Client
	broadcast   chan *Message
	redisClient *redis.Client
	users       []string
}

func NewHub(rdb *redis.Client) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan *Message),
		redisClient: rdb,
		users:       []string{},
	}
}

func (h *Hub) Run(ctx context.Context) {
	sub := h.redisClient.Subscribe(ChatChannel)
	ch := sub.Channel()

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.users = h.getUserList()
			h.broadcastUsersList()

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.conn.Close()
				h.users = h.getUserList()
				h.broadcastUsersList()
			}

		case msg := <-h.broadcast:
			data, _ := json.Marshal(msg)
			h.redisClient.Publish(ChatChannel, data)

		case redisMsg := <-ch:
			var m Message
			if err := json.Unmarshal([]byte(redisMsg.Payload), &m); err != nil {
				log.Println("Error decoding message:", err)
				continue
			}
			for c := range h.clients {
				c.send <- &m
			}
		}
	}
}

func (h *Hub) broadcastUsersList() {
	msg := &Message{
		Type:  "users",
		Users: h.getUserList(),
	}
	for c := range h.clients {
		c.send <- msg
	}
}

func (h *Hub) getUserList() []string {
	list := []string{}
	for c := range h.clients {
		list = append(list, c.name)
	}
	return list
}
