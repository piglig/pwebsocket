package ws

import "encoding/json"

type EventType string

const (
	SingleChatEvent EventType = "single_chat"
	GroupChatEvent  EventType = "group_chat"
	BroadChatEvent  EventType = "broadcast_chat"
)

type Event struct {
	Type    EventType
	*Client `json:"-"`
	Data    json.RawMessage
}

type EventHandler func(eventType EventType, client *Client, rawMsg json.RawMessage)

func isValidEventType(eventType EventType) bool {
	switch eventType {
	case SingleChatEvent:
	case GroupChatEvent:
	case BroadChatEvent:
	default:
		return false
	}
	return true
}
