package ws

import (
	"encoding/json"
	"reflect"
)

type Event interface {
	eventType() string
}

type RoomChanged struct{}

func (e *RoomChanged) eventType() string { return "RoomChanged" }

func Marshal(event Event) ([]byte, error) {
	obj := struct {
		Data      Event `json:",omitempty"`
		EventType string
	}{
		EventType: event.eventType(),
	}
	val := reflect.Indirect(reflect.ValueOf(event))
	if val.NumField() != 0 {
		obj.Data = event
	}
	return json.Marshal(&obj)
}
