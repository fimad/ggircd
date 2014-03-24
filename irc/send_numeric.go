package irc

import (
  "reflect"
)

// sendNumeric sends a numeric response to the given sink. If the sink is of
// type User
func sendNumeric(state State, sink Sink, msg Message, extra ...string) {
  params := make([]string, 0, len(extra)+1)

  // Attempt to add the nick name of the current sink to the error message.
  sinkType := reflect.TypeOf(sink)
  switch sinkType {
  case reflect.TypeOf(User{}):
    params = append(params, sink.(User).Nick)
  case reflect.TypeOf(&User{}):
    params = append(params, sink.(*User).Nick)
  }

  params = append(params, extra...)
  sink.Send(msg.WithPrefix(state.GetConfig().Name).WithParams(params...))
}

// sendNumericTrailing sends a numeric response to the given client.
func sendNumericTrailing(state State, sink Sink, msg Message, trailing string, extra ...string) {
  sendNumeric(state, sink, msg.WithTrailing(trailing), extra...)
}
