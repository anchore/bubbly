package bubbly

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wagoodman/go-partybus"
)

var (
	_ EventHandler = (*EventDispatcher)(nil)
	_ interface {
		EventHandler
		MessageListener
		HandleWaiter
	} = (*HandlerCollection)(nil)
)

type EventHandlerFn func(partybus.Event) []tea.Model

type EventHandler interface {
	partybus.Responder
	Handle(partybus.Event) []tea.Model
}

type MessageListener interface {
	OnMessage(tea.Msg)
}

type HandleWaiter interface {
	Wait()
}

type EventDispatcher struct {
	dispatch map[partybus.EventType]EventHandlerFn
	types    []partybus.EventType
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		dispatch: map[partybus.EventType]EventHandlerFn{},
	}
}

func (d *EventDispatcher) AddHandlers(handlers map[partybus.EventType]EventHandlerFn) {
	for t, h := range handlers {
		d.AddHandler(t, h)
	}
}

func (d *EventDispatcher) AddHandler(t partybus.EventType, fn EventHandlerFn) {
	d.dispatch[t] = fn
	d.types = append(d.types, t)
}

func (d EventDispatcher) RespondsTo() []partybus.EventType {
	return d.types
}

func (d EventDispatcher) Handle(e partybus.Event) []tea.Model {
	if fn, ok := d.dispatch[e.Type]; ok {
		return fn(e)
	}
	return nil
}

type HandlerCollection struct {
	handlers []EventHandler
}

func NewHandlerCollection(handlers ...EventHandler) *HandlerCollection {
	return &HandlerCollection{
		handlers: handlers,
	}
}

func (h *HandlerCollection) Append(handlers ...EventHandler) {
	h.handlers = append(h.handlers, handlers...)
}

func (h HandlerCollection) RespondsTo() []partybus.EventType {
	var ret []partybus.EventType
	for _, handler := range h.handlers {
		ret = append(ret, handler.RespondsTo()...)
	}
	return ret
}

func (h HandlerCollection) Handle(event partybus.Event) []tea.Model {
	var ret []tea.Model
	for _, handler := range h.handlers {
		ret = append(ret, handler.Handle(event)...)
	}
	return ret
}

func (h HandlerCollection) OnMessage(msg tea.Msg) {
	for _, handler := range h.handlers {
		if listener, ok := handler.(MessageListener); ok {
			listener.OnMessage(msg)
		}
	}
}

func (h HandlerCollection) Wait() {
	for _, handler := range h.handlers {
		if listener, ok := handler.(HandleWaiter); ok {
			listener.Wait()
		}
	}
}
