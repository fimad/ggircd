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
		got := state.GetChannel(channel).Mode
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
		if !state.GetChannel(channel).Ops[state.GetUser(nick)] {
			return fmt.Errorf("user %q should be op on %q but isn't", nick, channel)
		}
		return nil
	}
}

func assertChannelVoice(channel, nick string) assert {
	return func(state *mockState) error {
		if !state.GetChannel(channel).Voices[state.GetUser(nick)] {
			return fmt.Errorf("user %q should be voice on %q but doesn't", nick, channel)
		}
		return nil
	}
}
