package irc

type handler interface {
	// handle takes a connection and a messages, updates any relevant IRC state,
	// sends messages, etc. and returns a new handler that should be used from the
	// connection going forward.
	handle(connection, message) handler

	// closed is called when the connection has detected that the remote end point
	// is unreachable. This method is not called if the handler kills the
	// connection.
	closed(connection)
}

// nullHandler is a handler that ignores all messages.
type nullHandler struct{}

func (_ nullHandler) handle(_ connection, _ message) handler {
	return nullHandler{}
}

func (_ nullHandler) closed(c connection) {
	c.kill()
}

// sliceHandler is an implementation of handler that stores all received
// messages in a slice.
type sliceHandler []message

func (h *sliceHandler) handle(_ connection, msg message) handler {
	*h = append(*h, msg)
	return h
}

func (_ sliceHandler) closed(c connection) {
	c.kill()
}
