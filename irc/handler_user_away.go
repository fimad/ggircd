package irc

const (
	replyYouAreAway   = "You have been marked as being away"
	replyYouAreUnaway = "You are no longer marked as being away"
)

func (h *UserHandler) handleCmdAway(state State, user *User, conn Connection, msg Message) Handler {
	user.AwayMessage = msg.Trailing
	if len(msg.Trailing) == 0 {
		delete(user.Mode, UserModeAway)
		sendNumericTrailing(state, user, ReplyUnaway, replyYouAreUnaway)
	} else {
		user.Mode[UserModeAway] = true
		sendNumericTrailing(state, user, ReplyNowAway, replyYouAreAway)
	}
	return h
}
