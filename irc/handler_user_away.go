package irc

const (
	replyYouAreAway   = "You have been marked as being away"
	replyYouAreUnaway = "You are no longer marked as being away"
)

func (h *userHandler) handleCmdAway(state state, user *user, conn connection, msg message) handler {
	user.awayMessage = msg.trailing
	if len(msg.trailing) == 0 {
		delete(user.mode, userModeAway)
		sendNumericTrailing(state, user, replyUnaway, replyYouAreUnaway)
	} else {
		user.mode[userModeAway] = true
		sendNumericTrailing(state, user, replyNowAway, replyYouAreAway)
	}
	return h
}
