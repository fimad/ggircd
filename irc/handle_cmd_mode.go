package irc

import (
  "log"
  "strconv"
)

func (d *Dispatcher) handleCmdMode(msg Message, client *Client) {
  if client == nil {
    Todo("handle nil clients")
    return
  }

  if len(msg.Params) < 2 {
    client.Relay.Inbox <- ErrorNeedMoreParams
    return
  }

  if msg.Params[0][0] == '#' || msg.Params[0][0] == '&' {
    d.handleCmdModeChannel(msg, client)
  } else {
    d.handleCmdModeUser(msg, client)
  }
}

func (d *Dispatcher) handleCmdModeChannel(msg Message, client *Client) {
  channel := d.channels[msg.Params[0]]
  if channel == nil {
    msg.Relay.Inbox <- ErrorNoSuchChannel
    return
  }

  if !channel.Ops[client.ID] {
    msg.Relay.Inbox <- ErrorChanOPrivsNeeded.WithParams(channel.Name, "Not op")
    return
  }

  modes := msg.Params[1]

  affinity := true
  if modes[0] == '-' || modes[0] == '+' {
    affinity = modes[0] == '+'
    modes = modes[1:]
  } else {
    channel.Mode = make(Mode)
  }

  params := msg.Params[2:]
  numParams := 0

  // Ensures that all of the flags are valid and count the number of params
  // needed.
  for _, f := range modes[1:] {
    flag := ModeFlag(f)
    if flag == ChannelModeOp || flag == ChannelModeBanMask ||
      flag == ChannelModeUserLimit || flag == ChannelModeKey {
      numParams++
    } else if !ValidChannelModes[flag] {
      msg.Relay.Inbox <- ErrorUnknownMode.WithParams(string(f), "unknown mode")
      return
    }
  }

  if len(params) < numParams {
    msg.Relay.Inbox <- ErrorNeedMoreParams
    return
  }

  curParam := 0

  // Actually attempt to update the modes.
  for _, f := range modes {
    flag := ModeFlag(f)
    switch flag {
    case ChannelModeOp:
      nick := params[curParam]
      if cid, ok := d.nicks[nick]; ok {
        channel.Ops[cid] = affinity
      } else {
        msg.Relay.Inbox <- ErrorNoSuchNick.WithParams(nick)
      }
      curParam++

    case ChannelModeUserLimit:
      channel.Mode[flag] = affinity
      if !affinity {
        break
      }

      limit, err := strconv.Atoi(params[curParam])
      if err != nil {
        msg.Relay.Inbox <- ErrorUnknownMode.WithParams(params[curParam], "NaN")
      }
      channel.Limit = limit
      curParam++

    case ChannelModeBanMask:
      Todo("Handle ban masks")

    case ChannelModeVoice:
      nick := params[curParam]
      if cid, ok := d.nicks[nick]; ok {
        channel.Voice[cid] = affinity
      } else {
        msg.Relay.Inbox <- ErrorNoSuchNick.WithParams(nick)
      }
      curParam++

    case ChannelModeKey:
      channel.Mode[flag] = affinity
      if !affinity {
        break
      }
      channel.Key = params[curParam]
      curParam++

    case ChannelModePrivate:
      fallthrough
    case ChannelModeSecret:
      fallthrough
    case ChannelModeInvite:
      fallthrough
    case ChannelModeTopicOp:
      log.Printf("set affin...%q", affinity)
      fallthrough
    case ChannelModeNoOutside:
      fallthrough
    case ChannelModeModerated:
      channel.Mode[flag] = affinity
    }
  }

  // Broadcast the mode change, are we suppose to do this?
  msg.Prefix = client.Prefix()
  for cid := range channel.Clients {
    d.clients[cid].Relay.Inbox <- msg
  }
}

func (d *Dispatcher) handleCmdModeUser(msg Message, client *Client) {
  Todo("Implement user modes")
}
