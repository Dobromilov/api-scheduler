package handlers

import (
	"Scheduler-api/internal/store"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSClient struct {
	conn *websocket.Conn
	send chan []byte
	mu   sync.Mutex
}

func (c *WSClient) WriteMessage(messageType int, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(messageType, data)
}

type Hub struct {
	simID     string
	clients   map[*WSClient]bool
	mu        sync.RWMutex
	register  chan *WSClient
	unreg     chan *WSClient
	broadcast chan []byte
}

func NewHub(simID string) *Hub {
	return &Hub{
		simID:     simID,
		clients:   make(map[*WSClient]bool),
		register:  make(chan *WSClient),
		unreg:     make(chan *WSClient),
		broadcast: make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unreg:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Broadcast(data []byte) {
	select {
	case h.broadcast <- data:
	default:
	}
}

type GlobalHubManager struct {
	hubs map[string]*Hub
	mu   sync.RWMutex
}

var GlobalHubs = &GlobalHubManager{
	hubs: make(map[string]*Hub),
}

func (g *GlobalHubManager) GetOrCreateHub(simID string) *Hub {
	g.mu.Lock()
	defer g.mu.Unlock()

	if hub, ok := g.hubs[simID]; ok {
		return hub
	}

	hub := NewHub(simID)
	g.hubs[simID] = hub
	go hub.Run()
	return hub
}

func (g *GlobalHubManager) RemoveHub(simID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.hubs, simID)
}

func (g *GlobalHubManager) GetHub(simID string) (*Hub, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	hub, ok := g.hubs[simID]
	return hub, ok
}

func HandleSimulationWS(w http.ResponseWriter, r *http.Request) {
	simID := r.PathValue("id")
	if simID == "" {
		ErrorResponse(w, http.StatusBadRequest, "simulation ID is required")
		return
	}

	_, ok := store.GetStore().Get(simID)
	if !ok {
		ErrorResponse(w, http.StatusNotFound, "simulation not found")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	hub := GlobalHubs.GetOrCreateHub(simID)
	client := &WSClient{
		conn: conn,
		send: make(chan []byte, 256),
	}

	hub.register <- client

	go func() {
		defer func() {
			hub.unreg <- client
			conn.Close()
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()

	go func() {
		defer func() {
			hub.unreg <- client
			conn.Close()
		}()

		for message := range client.send {
			if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}()

	welcome, _ := json.Marshal(map[string]string{
		"type":    "connected",
		"message": "WebSocket connected for simulation " + simID,
	})
	client.WriteMessage(websocket.TextMessage, welcome)
}

func BroadcastProgress(simID string, data map[string]interface{}) {
	hub, ok := GlobalHubs.GetHub(simID)
	if !ok {
		return
	}

	msg, err := json.Marshal(map[string]interface{}{
		"type": "progress",
		"data": data,
	})
	if err != nil {
		return
	}

	hub.Broadcast(msg)
}
