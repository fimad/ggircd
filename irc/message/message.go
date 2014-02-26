package message

type Message struct {
  Prefix string
  Command string
  Params []string

  // The connection that this message originated from
  // Sender ...
}
