package push

import "github.com/Jeffail/gabs"

type ClientArray map[*Client]bool

type MessageContent *gabs.Container

type PushMessage struct {
	receiverId string
	message    MessageContent
	sender     *Client
}

type Hub struct {
	register   chan *Client
	unregister chan *Client
	send       chan *PushMessage
	clients    map[string]ClientArray
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		send:       make(chan *PushMessage),
		clients:    make(map[string]ClientArray),
	}
}

func (h *Hub) registerClient(c *Client) {
	if arr, ok := h.clients[c.id]; ok {
		arr[c] = true
	} else {
		arr := make(map[*Client]bool)
		arr[c] = true
		h.clients[c.id] = arr
	}
}

func (h *Hub) unregisterClient(c *Client) {
	if _, ok := h.clients[c.id]; ok {
		delete(h.clients[c.id], c)
	}
}

func (h *Hub) sendMessage(m *PushMessage) {
	if clients, ok := h.clients[m.receiverId]; ok && len(clients) > 0 {
		for client := range clients {
			client.send <- m.message
		}
	} else {
		// respond: no receiver!
	}
}

func (h *Hub) Run() {
	go func() {
		for {
			select {
			case c := <-h.register:
				h.registerClient(c)
			case c := <-h.unregister:
				h.unregisterClient(c)
			case m := <-h.send:
				h.sendMessage(m)
			}
		}
	}()
}