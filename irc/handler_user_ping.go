package irc

func (h *UserHandler) handleCmdPing(state State, user *User, conn Connection, msg Message) Handler {
  name := state.GetConfig().Name
  conn.Send(CmdPong.WithPrefix(name).WithParams(name).WithTrailing(name))
  return h
}
