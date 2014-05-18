package irc

import (
	"time"
)

type commandMap map[string]func(state, *user, connection, message) handler

// userHandler is a handler that handles messages coming from a user connection
// that has successfully associated with the client.
type userHandler struct {
	state      chan state
	nick       string
	commands   commandMap
	gotPong    chan struct{}
	killPing   chan struct{}
	shouldKill chan struct{}
}

func newUserHandler(state chan state, nick string) handler {
	handler := &userHandler{
		state:      state,
		nick:       nick,
		gotPong:    make(chan struct{}, 1),
		killPing:   make(chan struct{}, 1),
		shouldKill: make(chan struct{}, 1),
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
		cmdPong.command:    handler.handleCmdPong,
		cmdPrivMsg.command: handler.handleCmdPrivMsg,
		cmdQuit.command:    handler.handleCmdQuit,
		cmdTopic.command:   handler.handleCmdTopic,
		cmdWho.command:     handler.handleCmdWho,
	}
	go handler.pingLoop()
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

	select {
	case <-h.shouldKill:
		conn.kill()
		return h
	default:
	}

	command := h.commands[msg.command]
	if command == nil {
		return h
	}

	user := state.getUser(h.nick)
	newHandler := command(state, user, conn, msg)
	h.nick = user.nick

	if newHandler != h {
		h.killPing <- struct{}{}
	}
	return newHandler
}

func (h *userHandler) pingLoop() {
	state := <-h.state
	config := state.getConfig()
	h.state <- state

	pingDuration := time.Duration(config.PingFrequency) * time.Second
	pongDuration := time.Duration(config.PongMaxLatency) * time.Second
	pingTicker := time.NewTicker(pingDuration)
	var pongTimer <-chan time.Time

	alive := true
	for alive {
		select {
		case <-h.killPing:
			alive = false

		case <-pongTimer:
			h.shouldKill <- struct{}{}
			alive = false

		case <-h.gotPong:
			pongTimer = nil

		case <-pingTicker.C:
			state := <-h.state
			user := state.getUser(h.nick)
			user.send(cmdPing.withTrailing(config.Name))
			pongTimer = time.After(pongDuration)
			h.state <- state
		}
	}

	logf(debug, "Closing ping loop.")
	pingTicker.Stop()
}
