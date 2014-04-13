package irc

var (
	CmdInvite  = Message{Command: "INVITE"}
	CmdJoin    = Message{Command: "JOIN"}
	CmdList    = Message{Command: "LIST"}
	CmdMode    = Message{Command: "MODE"}
	CmdMotd    = Message{Command: "MOTD"}
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
