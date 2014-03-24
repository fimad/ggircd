package irc

type Handler interface {
	// Handle takes a connection and a messages, updates any relevant IRC state,
	// sends messages, etc. and returns a new handler that should be used from the
	// connection going forward.
	Handle(Connection, Message) Handler

	// Closed is called when the connection has detected that the remote end point
	// is unreachable. This method is not called if the handler kills the
	// connection.
	Closed(Connection)
}

// NullHandler is a handler that ignores all messages.
type NullHandler struct{}

func (_ NullHandler) Handle(_ Connection, _ Message) Handler {
	return NullHandler{}
}

func (_ NullHandler) Closed(c Connection) {
	c.Kill()
}

// SliceHandler is an implementation of Handler that stores all received
// messages in a slice.
type SliceHandler []Message

func (h *SliceHandler) Handle(_ Connection, msg Message) Handler {
	*h = append(*h, msg)
	return h
}

func (_ SliceHandler) Closed(c Connection) {
	c.Kill()
}
