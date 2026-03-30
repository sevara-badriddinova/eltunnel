package tunnel

import (
	"eltunnel/server/protocol"
	"sync"

	"github.com/gorilla/websocket"
)

func NewManager() *Manager {
	return &Manager{
		clients:         make(map[string]*websocket.Conn),
		pendingRequests: make(map[string]chan protocol.HttpResponsePayload),
	}
}

type Manager struct {
	clients         map[string]*websocket.Conn
	mu              sync.RWMutex
	pendingRequests map[string]chan protocol.HttpResponsePayload
}

func (m *Manager) Register(subdomain string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[subdomain] = conn
}

func (m *Manager) Remove(subdomain string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, subdomain)
}

func (m *Manager) GetClient(subdomain string) (*websocket.Conn, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conn, exists := m.clients[subdomain]
	return conn, exists
}

func (m *Manager) AddPendingRequest(requestID string, ch chan protocol.HttpResponsePayload) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pendingRequests[requestID] = ch
}

func (m *Manager) ResolvePendingRequest(requestID string, response protocol.HttpResponsePayload) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ch, exists := m.pendingRequests[requestID]
	if exists {
		ch <- response
	}
	delete(m.pendingRequests, requestID)
}
