package irc

var (
	cmdAway    = message{command: "AWAY"}
	cmdInvite  = message{command: "INVITE"}
	cmdJoin    = message{command: "JOIN"}
	cmdKick    = message{command: "KICK"}
	cmdList    = message{command: "LIST"}
	cmdMode    = message{command: "MODE"}
	cmdMotd    = message{command: "MOTD"}
	cmdNames   = message{command: "NAMES"}
	cmdNick    = message{command: "NICK"}
	cmdNotice  = message{command: "NOTICE"}
	cmdPart    = message{command: "PART"}
	cmdPass    = message{command: "PASS"}
	cmdPing    = message{command: "PING"}
	cmdPong    = message{command: "PONG"}
	cmdPrivMsg = message{command: "PRIVMSG"}
	cmdQuit    = message{command: "QUIT"}
	cmdServer  = message{command: "SERVER"}
	cmdTopic   = message{command: "TOPIC"}
	cmdUser    = message{command: "USER"}
	cmdWho     = message{command: "WHO"}
)
