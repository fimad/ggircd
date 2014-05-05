package irc

const errInviteNotOnChannel = "You're not on that channel"
const errInviteChanOPrivs = "You're not channel operator"
const errInviteUserOnChannel = "is already on channel"

func (h *userHandler) handleCmdInvite(state state, user *user, conn connection, msg message) handler {
	if len(msg.params) < 2 {
		sendNumeric(state, user, errorNeedMoreParams)
		return h
	}

	nickName := msg.params[0]
	channelName := msg.params[1]
	nick := state.getUser(nickName)
	channel := state.getChannel(channelName)

	if nick == nil {
		sendNumeric(state, user, errorNoSuchNick)
		return h
	}

	// If the channel does not exist, treat it as a success and just send an
	// invite.
	if channel == nil {
		sendNumeric(state, user, replyInviting, channelName, nickName)
		nick.send(msg.withPrefix(user.prefix()))
		return h
	}

	// Don't invite users if they are already on the channel.
	if channel.users[nick] {
		sendNumericTrailing(state, user,
			errorUserOnChannel, errInviteUserOnChannel, nickName, channelName)
		return h
	}

	// If the channel does exist, then the user performing the invite must be on
	// that channel.
	if !channel.users[user] {
		sendNumericTrailing(state, user,
			errorNotOnChannel, errInviteNotOnChannel, channelName)
		return h
	}

	// If the channel is set as invite only, then only channel operators are
	// allowed to invite.
	if channel.mode[channelModeInvite] && !channel.ops[user] {
		sendNumericTrailing(state, user,
			errorChanOPrivsNeeded, errInviteChanOPrivs, channelName)
		return h
	}

	channel.invite(nick)
	sendNumeric(state, user, replyInviting, channelName, nickName)
	nick.send(msg.withPrefix(user.prefix()))
	return h
}
