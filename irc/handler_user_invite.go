package irc

const errInviteNotOnChannel = "You're not on that channel"
const errInviteChanOPrivs = "You're not channel operator"
const errInviteUserOnChannel = "is already on channel"

func (h *UserHandler) handleCmdInvite(state State, user *User, conn Connection, msg Message) Handler {
	if len(msg.Params) < 2 {
		sendNumeric(state, user, ErrorNeedMoreParams)
		return h
	}

	nickName := msg.Params[0]
	channelName := msg.Params[1]
	nick := state.GetUser(nickName)
	channel := state.GetChannel(channelName)

	if nick == nil {
		sendNumeric(state, user, ErrorNoSuchNick)
		return h
	}

	// If the channel does not exist, treat it as a success and just send an
	// invite.
	if channel == nil {
		sendNumeric(state, user, ReplyInviting, channelName, nickName)
		nick.Send(msg.WithPrefix(user.Prefix()))
		return h
	}

	// Don't invite users if they are already on the channel.
	if channel.Users[nick] {
		sendNumericTrailing(state, user,
			ErrorUserOnChannel, errInviteUserOnChannel, nickName, channelName)
		return h
	}

	// If the channel does exist, then the user performing the invite must be on
	// that channel.
	if !channel.Users[user] {
		sendNumericTrailing(state, user,
			ErrorNotOnChannel, errInviteNotOnChannel, channelName)
		return h
	}

	// If the channel is set as invite only, then only channel operators are
	// allowed to invite.
	if channel.Mode[ChannelModeInvite] && !channel.Ops[user] {
		sendNumericTrailing(state, user,
			ErrorChanOPrivsNeeded, errInviteChanOPrivs, channelName)
		return h
	}

	channel.Invite(nick)
	sendNumeric(state, user, ReplyInviting, channelName, nickName)
	nick.Send(msg.WithPrefix(user.Prefix()))
	return h
}
