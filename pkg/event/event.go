package event

import (
	"fmt"
	"sync"
)

type EventId string

type Event struct {
	Id   EventId
	Data interface{}
	lock sync.Mutex // just placed for now, i don't know why i place it here, seems useless
}

type HandlerFunc func(event *Event)

type Handler struct {
	EventId EventId
	Name    string
	Handler HandlerFunc
}

type Manager struct {
	handlers         map[EventId]map[string]*Handler
	eventWorkerIdMap map[EventId]int
	pendingEvents    []*Event
	workerQueue      []chan func()
	currentWorkerId  int
	stopSig          chan int
	queueSize        int
	workerSize       int
	dispatching      bool
	lock             sync.RWMutex
}

func NewManger(queueSize int, workerSize int) *Manager {
	manager := &Manager{
		handlers:         make(map[EventId]map[string]*Handler),
		eventWorkerIdMap: make(map[EventId]int),
		workerQueue:      make([]chan func(), workerSize),
		currentWorkerId:  0,
		stopSig:          make(chan int, workerSize),
		queueSize:        queueSize,
		workerSize:       workerSize,
		lock:             sync.RWMutex{},
		dispatching:      false,
		pendingEvents:    make([]*Event, 0),
	}
	for i := 0; i < workerSize; i++ {
		queue := make(chan func(), queueSize)
		manager.workerQueue[i] = queue
		go func() {
			for {
				select {
				case <-manager.stopSig:
					return
				case f := <-queue:
					f()
				}
			}
		}()
	}
	return manager
}

// Start for starting to dispatching events
func (h *Manager) Start() {
	h.dispatching = true
	for _, event := range h.pendingEvents {
		h.Call(event)
	}
}

func (h *Manager) Stop() {
	for i := 0; i < h.workerSize; i++ {
		h.stopSig <- 0
	}
	h.dispatching = false
}

func (h *Manager) Register(handler *Handler) {
	h.lock.Lock()
	m, ok := h.handlers[handler.EventId]
	// if not found, crate new handler map and assign a worker to this event id
	if !ok {
		m = make(map[string]*Handler)
		h.handlers[handler.EventId] = m
		// assign a worker to this event id
		h.eventWorkerIdMap[handler.EventId] = h.currentWorkerId
		h.currentWorkerId = (h.currentWorkerId + 1) % h.workerSize
	}
	if _, ok := m[handler.Name]; ok {
		fmt.Printf("handler %s already registered, old handler is overwrittened\n", handler.Name)
	}
	m[handler.Name] = handler
	h.lock.Unlock()
}

func (h *Manager) RegisterA(id EventId, name string, handler HandlerFunc) {
	h.Register(&Handler{
		EventId: id,
		Name:    name,
		Handler: handler,
	})
}

func (h *Manager) UnregisterAll() {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.handlers = make(map[EventId]map[string]*Handler)
	h.currentWorkerId = 0
	h.eventWorkerIdMap = make(map[EventId]int)
}

func (h *Manager) Unregister(name string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	for _, m := range h.handlers {
		if _, ok := m[name]; ok {
			delete(m, name)
		}
	}
}

func (h *Manager) Call(event *Event) {

	// if not dispatching, put this event to pending events
	h.lock.Lock()
	if !h.dispatching {
		h.pendingEvents = append(h.pendingEvents, event)
		h.lock.Unlock()
		return
	}

	handlers, ok := h.handlers[event.Id]
	if !ok {
		h.lock.Unlock()
		return
	}
	workerId, ok := h.eventWorkerIdMap[event.Id]
	if !ok {
		// event id don't have a worker id, ignore
		// maybe because this event id has no handler
		h.lock.Unlock()
		return
	}

	h.lock.Unlock()
	for _, eh := range handlers {
		eventHandler := eh
		h.workerQueue[workerId] <- func() {
			event.lock.Lock()
			eventHandler.Handler(event)
			event.lock.Unlock()
		}
	}
}

func (h *Manager) CallA(id EventId, data interface{}) {
	h.Call(&Event{
		Id:   id,
		Data: data,
	})
}
