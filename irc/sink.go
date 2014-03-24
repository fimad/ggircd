package irc

// Sink is a abstraction over network connections that is able to take a Message
// and forward it to the appropriated connection.
type Sink interface {
	Send(Message)
}

// NullSink is an implementation of sink that drops all messages on the floor.
type NullSink struct{}

func (_ NullSink) Send(msg Message) {}

// SliceSink is an implementation of sink that stores all received messages in a
// slice.
type SliceSink []Message

func (s *SliceSink) Send(msg Message) {
	*s = append(*s, msg)
}
