package irc

import (
	"reflect"
)

// sendNumeric sends a numeric response to the given sink. If the sink is of
// type User
func sendNumeric(state state, sink sink, msg message, extra ...string) {
	params := make([]string, 0, len(extra)+1)

	// Attempt to add the nick name of the current sink to the error message.
	sinkType := reflect.TypeOf(sink)
	switch sinkType {
	case reflect.TypeOf(user{}):
		params = append(params, sink.(user).nick)
	case reflect.TypeOf(&user{}):
		params = append(params, sink.(*user).nick)
	}

	params = append(params, extra...)
	sink.send(msg.withPrefix(state.getConfig().Name).withParams(params...))
}

// sendNumericTrailing sends a numeric response to the given client.
func sendNumericTrailing(state state, sink sink, msg message, trailing string, extra ...string) {
	sendNumeric(state, sink, msg.withTrailing(trailing), extra...)
}
