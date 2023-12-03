package schema

import (
	"encoding/json"
	"reflect"
)

type Event interface {
	EventType() string
}

type RoomChanged struct{}

func (e *RoomChanged) EventType() string { return "RoomChanged" }

type GameStarted struct{}

func (e *GameStarted) EventType() string { return "GameStarted" }

type GameChanged struct{}

func (e *GameChanged) EventType() string { return "GameChanged" }

type Heartbeat struct{}

func (e *Heartbeat) EventType() string { return "Heartbeat" }

func MarshalEvent(event Event) ([]byte, error) {
	obj := struct {
		Data      Event `json:",omitempty"`
		EventType string
	}{
		EventType: event.EventType(),
	}
	val := reflect.Indirect(reflect.ValueOf(event))
	if val.NumField() != 0 {
		obj.Data = event
	}
	return json.Marshal(&obj)
}
