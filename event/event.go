package event

import (
	"AynaLivePlayer/logger"
	"github.com/sirupsen/logrus"
	"sync"
)

type EventId string

const MODULE_HANDLER = "EventHandler"

var eventLogger = logger.Logger.WithFields(logrus.Fields{
	"Module": MODULE_HANDLER,
})

type Event struct {
	Id        EventId
	Cancelled bool
	Data      interface{}
}

type EventHandlerFunc func(event *Event)

type EventHandler struct {
	EventId EventId
	Name    string
	Handler EventHandlerFunc
}

type Handler struct {
	handlers map[string]*EventHandler
	lock     sync.RWMutex
}

func NewHandler() *Handler {
	return &Handler{
		handlers: make(map[string]*EventHandler),
	}
}

func (h *Handler) Register(handler *EventHandler) {
	h.lock.Lock()
	defer h.lock.Unlock()
	eventLogger.Tracef("register new handler id=%s,name=%s", handler.EventId, handler.Name)
	h.handlers[handler.Name] = handler
}

func (h *Handler) RegisterA(id EventId, name string, handler EventHandlerFunc) {
	h.Register(&EventHandler{
		EventId: id,
		Name:    name,
		Handler: handler,
	})
}

func (h *Handler) UnregisterAll() {
	h.lock.Lock()
	defer h.lock.Unlock()
	eventLogger.Trace("clear all handler")
	h.handlers = make(map[string]*EventHandler)
}

func (h *Handler) Unregister(name string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	eventLogger.Tracef("unregister handler name=%s", name)
	delete(h.handlers, name)
}

func (h *Handler) Call(event *Event) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	for _, eh := range h.handlers {
		if eh.EventId == event.Id {
			eventLogger.Tracef("handler name=%s called by event_id = %s", event.Id, eh.Name)
			// todo: @3
			go eh.Handler(event)
		}
	}
}

func (h *Handler) CallA(id EventId, data interface{}) {
	h.Call(&Event{
		Id:   id,
		Data: data,
	})
}
