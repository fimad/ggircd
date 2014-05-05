package irc

const (
	replyYouAreAway   = "You are no longer marked as being away"
	replyYouAreUnaway = "You have been marked as being away"
)

func (h *UserHandler) handleCmdAway(state State, user *User, conn Connection, msg Message) Handler {
	user.AwayMessage = msg.Trailing
	if msg.Trailing == "" {
		delete(user.Mode, UserModeAway)
		sendNumeric(state, user, ReplyUnaway, replyYouAreUnaway)
	} else {
		user.Mode[UserModeAway] = true
		sendNumeric(state, user, ReplyNowAway, replyYouAreAway)
	}
	return h
}
