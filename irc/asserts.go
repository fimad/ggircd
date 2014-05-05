package irc

import (
	"fmt"
	"reflect"
)

// assert is a method that takes a mock state and produces an error if the state
// has an unexpected property.
type assert func(*mockState) error

func assertChannelMode(channel, modeLine string) assert {
	want := ParseMode(ChannelModes, modeLine)
	return func(state *mockState) error {
		ch := state.GetChannel(channel)
		if ch == nil {
			return fmt.Errorf("expected channel %q to exist, but does not", channel)
		}

		got := ch.Mode
		// delete all keys with false values from got
		for k, v := range got {
			if !v {
				delete(got, k)
			}
		}

		if !reflect.DeepEqual(got, want) {
			return fmt.Errorf("unexpected mode on %s =>\n\t\tgot %+v\n\t\twant %+v",
				channel, got, want)
		}
		return nil
	}
}

func assertChannelOp(channel, nick string) assert {
	return func(state *mockState) error {
		ch := state.GetChannel(channel)
		if ch == nil {
			return fmt.Errorf("expected channel %q to exist, but does not", channel)
		}

		user := state.GetUser(nick)
		if user == nil {
			return fmt.Errorf("expected user %q to exist, but does not", nick)
		}

		if !ch.Ops[user] {
			return fmt.Errorf("user %q should be op on %q but isn't", nick, channel)
		}
		return nil
	}
}

func assertChannelVoice(channel, nick string) assert {
	return func(state *mockState) error {
		ch := state.GetChannel(channel)
		if ch == nil {
			return fmt.Errorf("expected channel %q to exist, but does not", channel)
		}

		user := state.GetUser(nick)
		if user == nil {
			return fmt.Errorf("expected user %q to exist, but does not", nick)
		}

		if !ch.Voices[user] {
			return fmt.Errorf("user %q should be voice on %q but doesn't", nick, channel)
		}
		return nil
	}
}

func assertUserOnChannel(nick, channel string) assert {
	return func(state *mockState) error {
		ch := state.GetChannel(channel)
		if ch == nil {
			return fmt.Errorf("expected channel %q to exist, but does not", channel)
		}

		user := state.GetUser(nick)
		if user == nil {
			return fmt.Errorf("expected user %q to exist, but does not", nick)
		}

		if !ch.Users[user] {
			return fmt.Errorf("user %q should be on %q but isn't", nick, channel)
		}
		return nil
	}
}

func assertUserNotOnChannel(nick, channel string) assert {
	return func(state *mockState) error {
		ch := state.GetChannel(channel)
		if ch == nil {
			return fmt.Errorf("expected channel %q to exist, but does not", channel)
		}

		user := state.GetUser(nick)
		if user == nil {
			return fmt.Errorf("expected user %q to exist, but does not", nick)
		}

		if ch.Users[user] {
			return fmt.Errorf("user %q should not be on %q but is", nick, channel)
		}
		return nil
	}
}

func assertUserMode(nick, modeLine string) assert {
	want := ParseMode(UserModes, modeLine)
	return func(state *mockState) error {
		user := state.GetUser(nick)
		if user == nil {
			return fmt.Errorf("expected user %q to exist, but does not", nick)
		}

		got := user.Mode
		// delete all keys with false values from got
		for k, v := range got {
			if !v {
				delete(got, k)
			}
		}

		if !reflect.DeepEqual(got, want) {
			return fmt.Errorf("unexpected mode on %s =>\n\t\tgot %+v\n\t\twant %+v",
				nick, got, want)
		}
		return nil
	}
}
