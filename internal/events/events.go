package events

import "github.com/grafana/sobek"

type EventListener struct {
	EventType string

	List []func(sobek.Value) (sobek.Value, error)
}

func (e *EventListener) Add(fn func(sobek.Value) (sobek.Value, error)) {
	e.List = append(e.List, fn)
}

func (e *EventListener) All() []func(sobek.Value) (sobek.Value, error) {
	return e.List
}
