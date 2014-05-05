package irc

// sink is a abstraction over network connections that is able to take a message
// and forward it to the appropriated connection.
type sink interface {
	send(message)
}

// nullSink is an implementation of sink that drops all messages on the floor.
type nullSink struct{}

func (_ nullSink) send(msg message) {}

// sliceSink is an implementation of sink that stores all received messages in a
// slice.
type sliceSink []message

func (s *sliceSink) send(msg message) {
	*s = append(*s, msg)
}
