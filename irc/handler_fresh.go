package irc

// FreshHandler is a handler for a brand new connection that has not been
// registered yet.
type FreshHandler struct {
}

func NewFreshHandler() Handler {
  return &FreshHandler{}
}

func (h *FreshHandler) Handle(conn Connection, msg Message) Handler {
  return h
}
