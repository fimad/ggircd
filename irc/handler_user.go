package irc

import (
	"github.com/prometheus/client_golang/prometheus"
)

type commandMap map[string]func(state, *user, connection, message) handler

// userHandler is a handler that handles messages coming from a user connection
// that has successfully associated with the client.
type userHandler struct {
	state    chan state
	nick     string
	commands commandMap
}

func newUserHandler(state chan state, nick string) handler {
	handler := &userHandler{
		state: state,
		nick:  nick,
	}
	handler.commands = commandMap{
		cmdAway.command:    handler.handleCmdAway,
		cmdInvite.command:  handler.handleCmdInvite,
		cmdJoin.command:    handler.handleCmdJoin,
		cmdKick.command:    handler.handleCmdKick,
		cmdList.command:    handler.handleCmdList,
		cmdMode.command:    handler.handleCmdMode,
		cmdMotd.command:    handler.handleCmdMotd,
		cmdNames.command:   handler.handleCmdNames,
		cmdNick.command:    handler.handleCmdNick,
		cmdNotice.command:  handler.handleCmdNotice,
		cmdPart.command:    handler.handleCmdPart,
		cmdPing.command:    handler.handleCmdPing,
		cmdPrivMsg.command: handler.handleCmdPrivMsg,
		cmdQuit.command:    handler.handleCmdQuit,
		cmdTopic.command:   handler.handleCmdTopic,
		cmdWho.command:     handler.handleCmdWho,
	}
	return handler
}

func (h *userHandler) closed(conn connection) {
	state := <-h.state
	defer func() { h.state <- state }()

	state.removeUser(state.getUser(h.nick))
	conn.kill()
}

func (h *userHandler) handle(conn connection, msg message) handler {
	state := <-h.state
	defer func() { h.state <- state }()

	command_received.With(
		prometheus.Labels{
			"nick":    h.nick,
			"command": msg.command,
		},
	).Inc()

	command := h.commands[msg.command]
	if command == nil {
		return h
	}

	user := state.getUser(h.nick)
	newHandler := command(state, user, conn, msg)
	h.nick = user.nick
	return newHandler
}
