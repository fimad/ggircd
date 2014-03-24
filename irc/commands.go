package irc

var (
  CmdJoin    = Message{Command: "JOIN"}
  CmdMode    = Message{Command: "MODE"}
  CmdNames   = Message{Command: "NAMES"}
  CmdNick    = Message{Command: "NICK"}
  CmdPart    = Message{Command: "PART"}
  CmdPass    = Message{Command: "PASS"}
  CmdPing    = Message{Command: "PING"}
  CmdPong    = Message{Command: "PONG"}
  CmdPrivMsg = Message{Command: "PRIVMSG"}
  CmdQuit    = Message{Command: "QUIT"}
  CmdServer  = Message{Command: "SERVER"}
  CmdTopic   = Message{Command: "TOPIC"}
  CmdUser    = Message{Command: "USER"}
)
