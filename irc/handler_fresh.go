package irc

// FreshHandler is a handler for a brand new connection that has not been
// registered yet.
type FreshHandler struct {
  state chan State
}

func NewFreshHandler(state chan State) Handler {
  return &FreshHandler{state: state}
}

func (h *FreshHandler) Handle(conn Connection, msg Message) Handler {
  if msg.Command != CmdNick.Command {
    return h
  }
  return h.handleNick(conn, msg)
}

func (_ *FreshHandler) Closed(c Connection) {
  c.Kill()
}

func (h *FreshHandler) handleNick(conn Connection, msg Message) Handler {
  state := <-h.state
  defer func() { h.state <- state }()

  if len(msg.Params) < 1 {
    sendNumeric(state, conn, ErrorNoNicknameGiven)
    return h
  }
  nick := msg.Params[0]

  user := state.NewUser(nick)
  if user == nil {
    sendNumeric(state, conn, ErrorNicknameInUse)
    return h
  }
  user.Sink = conn

  return &freshUserHandler{state: h.state, user: user}
}

// freshUserHandler is a handler for a brand new connection that is in the
// process of registering and has succesfully set a nickname.
type freshUserHandler struct {
  user  *User
  state chan State
}

func (h *freshUserHandler) Handle(conn Connection, msg Message) Handler {
  if msg.Command != CmdUser.Command {
    return h
  }
  return h.handleUser(conn, msg)
}

func (_ *freshUserHandler) Closed(c Connection) {
  c.Kill()
}

func (h *freshUserHandler) handleUser(conn Connection, msg Message) Handler {
  state := <-h.state
  defer func() { h.state <- state }()

  if len(msg.Params) < 3 || msg.Trailing == "" {
    sendNumeric(state, h.user, ErrorNeedMoreParams)
    return h
  }

  h.user.User = msg.Params[0]
  h.user.Host = msg.Params[1]
  h.user.Server = msg.Params[2]
  h.user.RealName = msg.Trailing

  sendIntro(state, h.user)

  return NullHandler{}
}
