package eventbus

// BusBridge is a minimal interface for a websocket-like JSON connection.
// Implemented by e.g. *websocket.Conn from gorilla via wrappers:
type BusBridge interface {
	ReadJSON(v any) error
	WriteJSON(v any) error
	Close() error
}
