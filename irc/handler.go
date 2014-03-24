package irc

type Handler interface {
  Handle(Connection, Message) Handler
}

// NullHandler is a handler that ignores all messages.
type NullHandler struct {}

func (_ NullHandler) Handle(_ Connection, _ Message) Handler {
  return NullHandler{}
}

// SliceHandler is an implementation of Handler that stores all received
// messages in a slice.
type SliceHandler []Message

func (h *SliceHandler) Handle(_ Connection,  msg Message) Handler {
  *h = append(*h, msg)
  return h
}
