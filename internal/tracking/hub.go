package tracking

import "sync"

type Hub struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan Event]struct{}
}

func NewHub() *Hub {
	return &Hub{
		subscribers: make(map[string]map[chan Event]struct{}),
	}
}

func (h *Hub) Subscribe(loadID string) chan Event {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan Event, 8)
	if _, exists := h.subscribers[loadID]; !exists {
		h.subscribers[loadID] = make(map[chan Event]struct{})
	}

	h.subscribers[loadID][ch] = struct{}{}
	return ch
}

func (h *Hub) Unsubscribe(loadID string, ch chan Event) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if subscribers, exists := h.subscribers[loadID]; exists {
		delete(subscribers, ch)
		if len(subscribers) == 0 {
			delete(h.subscribers, loadID)
		}
	}

	close(ch)
}

func (h *Hub) Publish(loadID string, event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for subscriber := range h.subscribers[loadID] {
		select {
		case subscriber <- event:
		default:
		}
	}
}
