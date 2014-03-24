package irc

import (
  "testing"
)

func TestUserHandlerJoin(t *testing.T) {
  state := make(chan State, 1)
  testHandler(t, "UserHandler-JOIN", state, NewUserHandler(state, "nick"), []handlerTest{
    {
      desc: "succesful join",
      in:   []Message{CmdJoin.WithParams("#channel")},
      wantNicks: map[string]mockConnection{
        "nick": mockConnection{
          messages: []Message{
            ReplyNoTopic,
            ReplyNamReply,
            ReplyEndOfNames,
          },
        },
      },
      state: newMockState().withUser("nick"),
    },
    {
      desc: "succesful join with key",
      in:   []Message{CmdJoin.WithParams("#channel", "key")},
      wantNicks: map[string]mockConnection{
        "nick": mockConnection{
          messages: []Message{
            ReplyNoTopic,
            ReplyNamReply,
            ReplyEndOfNames,
          },
        },
      },
      state: newMockState().
        withUser("nick").
        withChannel("#channel", "k").
        withChannelKey("#channel", "key"),
    },
    {
      desc: "succesful join limited channel",
      in:   []Message{CmdJoin.WithParams("#channel")},
      wantNicks: map[string]mockConnection{
        "nick": mockConnection{
          messages: []Message{
            ReplyNoTopic,
            ReplyNamReply,
            ReplyEndOfNames,
          },
        },
      },
      state: newMockState().
        withUser("nick").
        withChannel("#channel", "l").
        withChannelLimit("#channel", 1),
    },
    {
      desc: "succesful join multiple",
      in:   []Message{CmdJoin.WithParams("#foo,#bar")},
      wantNicks: map[string]mockConnection{
        "nick": mockConnection{
          messages: []Message{
            ReplyNoTopic,
            ReplyNamReply,
            ReplyEndOfNames,
            ReplyNoTopic,
            ReplyNamReply,
            ReplyEndOfNames,
          },
        },
      },
      state: newMockState().withUser("nick"),
    },
    {
      desc: "no channel given",
      in:   []Message{CmdJoin},
      wantNicks: map[string]mockConnection{
        "nick": mockConnection{
          messages: []Message{ErrorNeedMoreParams},
        },
      },
      state: newMockState().withUser("nick"),
    },
    {
      desc: "bad channel name",
      in:   []Message{CmdJoin.WithParams("channel")},
      wantNicks: map[string]mockConnection{
        "nick": mockConnection{
          messages: []Message{ErrorNoSuchChannel},
        },
      },
      state: newMockState().withUser("nick"),
    },
    {
      desc: "invite only",
      in:   []Message{CmdJoin.WithParams("#channel")},
      wantNicks: map[string]mockConnection{
        "nick": mockConnection{
          messages: []Message{ErrorInviteOnlyChan},
        },
      },
      state: newMockState().withUser("nick").withChannel("#channel", "i"),
    },
    {
      desc: "wrong channel key",
      in:   []Message{CmdJoin.WithParams("#channel")},
      wantNicks: map[string]mockConnection{
        "nick": mockConnection{
          messages: []Message{ErrorBadChannelKey},
        },
      },
      state: newMockState().
        withUser("nick").
        withChannel("#channel", "k").
        withChannelKey("#channel", "key"),
    },
    {
      desc: "channel full",
      in:   []Message{CmdJoin.WithParams("#channel")},
      wantNicks: map[string]mockConnection{
        "nick": mockConnection{
          messages: []Message{ErrorChannelIsFull},
        },
      },
      state: newMockState().
        withUser("nick").
        withChannel("#channel", "l").
        withChannelLimit("#channel", 0),
    },
  })
}
