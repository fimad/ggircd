package irc

// freshHandler is a handler for a brand new connection that has not been
// registered yet.
type freshHandler struct {
	state chan state
}

func newFreshHandler(state chan state) handler {
	return &freshHandler{state: state}
}

func (h *freshHandler) handle(conn connection, msg message) handler {
	if msg.command != cmdNick.command {
		return h
	}
	return h.handleNick(conn, msg)
}

func (_ *freshHandler) closed(c connection) {
	c.kill()
}

func (h *freshHandler) handleNick(conn connection, msg message) handler {
	state := <-h.state
	defer func() { h.state <- state }()

	if len(msg.params) < 1 {
		sendNumeric(state, conn, errorNoNicknameGiven)
		return h
	}
	nick := msg.params[0]

	user := state.newUser(nick)
	if user == nil {
		sendNumeric(state, conn, errorNicknameInUse)
		return h
	}
	user.sink = conn

	return &freshUserHandler{state: h.state, user: user}
}

// freshUserHandler is a handler for a brand new connection that is in the
// process of registering and has successfully set a nickname.
type freshUserHandler struct {
	user  *user
	state chan state
}

func (h *freshUserHandler) handle(conn connection, msg message) handler {
	if msg.command != cmdUser.command {
		return h
	}
	return h.handleUser(conn, msg)
}

func (h *freshUserHandler) closed(c connection) {
	state := <-h.state
	defer func() { h.state <- state }()

	state.removeUser(h.user)
	c.kill()
}

func (h *freshUserHandler) handleUser(conn connection, msg message) handler {
	state := <-h.state
	defer func() { h.state <- state }()
	logf(warn, "USER!!! %+v", msg)

	if len(msg.params) < 3 || msg.trailing == "" {
		sendNumeric(state, h.user, errorNeedMoreParams)
		return h
	}

	h.user.user = msg.params[0]
	h.user.host = state.getConfig().SpoofHostName
	if h.user.host == "" {
		h.user.host = msg.params[1]
	}
	h.user.server = state.getConfig().Name
	h.user.realName = msg.trailing

	sendIntro(state, h.user)

	return newUserHandler(h.state, h.user.nick)
}
