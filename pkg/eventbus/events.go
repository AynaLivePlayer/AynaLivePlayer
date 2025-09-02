package eventbus

import (
	"encoding/json"
	"errors"
	"reflect"
)

var DefaultMapper = NewEventsMapper()

func UnmarshalEvent(data []byte) (*Event, error) {
	return DefaultMapper.UnmarshalEvent(data)
}

func UnmarshalEventData(eventId string, data []byte) (any, error) {
	return DefaultMapper.UnmarshalEventData(eventId, data)
}

type EventsMapper struct {
	Mapping map[string]any
}

func NewEventsMapper() *EventsMapper {
	return &EventsMapper{
		Mapping: make(map[string]any),
	}
}

type untypedEvent struct {
	Id      string
	Channel string
	EchoId  string
	Data    json.RawMessage
}

func (m *EventsMapper) UnmarshalEvent(data []byte) (*Event, error) {
	var val untypedEvent
	err := json.Unmarshal(data, &val)
	if err != nil {
		return nil, errors.New("failed to unmarshal event: " + err.Error())
	}
	actualEventData, err := m.UnmarshalEventData(val.Id, val.Data)
	if err != nil {
		return nil, errors.New("failed to unmarshal event data: " + err.Error())
	}
	return &Event{
		Id:      val.Id,
		Channel: val.Channel,
		EchoId:  val.EchoId,
		Data:    actualEventData,
	}, nil
}

func (m *EventsMapper) UnmarshalEventData(eventId string, data []byte) (any, error) {
	val, ok := m.Mapping[eventId]
	if !ok {
		return nil, errors.New("event id not found")
	}
	newVal := reflect.New(reflect.TypeOf(val))
	err := json.Unmarshal(data, newVal.Interface())
	if err != nil {
		return nil, err
	}
	return newVal.Elem().Interface(), nil
}
