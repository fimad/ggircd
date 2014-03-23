package irc

import (
  "strconv"
)

func (d *Dispatcher) handleCmdMode(msg Message, client *Client) {
  if client == nil {
    Todo("handle nil clients")
    return
  }

  if len(msg.Params) < 1 {
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
  channel := d.ChannelForName(msg.Params[0])
  if channel == nil {
    d.sendNumeric(client, ErrorNoSuchChannel, msg.Params[0])
    return
  }

  if len(msg.Params) == 1 {
    d.sendChannelMode(client, channel)
    return
  }

  if !channel.Ops[client.ID] {
    d.sendNumeric(client, ErrorChanOPrivsNeeded, channel.Name)
    return
  }

  modes := msg.Params[1]
  params := msg.Params[2:]
  numParams := 0

  // Ensures that all of the flags are valid and count the number of params
  // needed.
  for _, f := range modes {
    flag := ModeFlag(f)
    if flag == "-" || flag == "+" {
      continue
    }

    if flag == ChannelModeOp || flag == ChannelModeBanMask ||
      flag == ChannelModeUserLimit || flag == ChannelModeKey {
      numParams++
    } else if !ValidChannelModes[flag] {
      d.sendNumeric(client, ErrorUnknownMode, string(f))
      return
    }
  }

  if len(params) < numParams {
    d.sendNumeric(client, ErrorNeedMoreParams)
    return
  }

  curParam := 0

  // Actually attempt to update the modes.
  affinity := true
  for _, f := range modes {
    flag := ModeFlag(f)
    switch flag {
    case "+":
      affinity = true
    case "-":
      affinity = false
    case ChannelModeOp:
      nick := params[curParam]
      if nickClient, ok := d.ClientForNick(nick); ok {
        channel.Ops[nickClient.ID] = affinity
      } else {
        d.sendNumeric(client, ErrorNoSuchNick, nick)
      }
      curParam++

    case ChannelModeUserLimit:
      channel.Mode[flag] = affinity
      if !affinity {
        break
      }

      limit, err := strconv.Atoi(params[curParam])
      if err != nil {
        d.sendNumericTrailing(client, ErrorUnknownMode, "NaN", params[curParam])
      }
      channel.Limit = limit
      curParam++

    case ChannelModeBanMask:
      Todo("Handle ban masks")

    case ChannelModeVoice:
      nick := params[curParam]
      if client, ok := d.ClientForNick(nick); ok {
        channel.Voice[client.ID] = affinity
      } else {
        d.sendNumeric(client, ErrorNoSuchNick, nick)
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
      fallthrough
    case ChannelModeNoOutside:
      fallthrough
    case ChannelModeModerated:
      channel.Mode[flag] = affinity
    }
  }

  // Broadcast the mode change, are we suppose to do this?
  msg.Prefix = client.Prefix()
  d.SendToChannel(channel, msg)
}

func (d *Dispatcher) handleCmdModeUser(msg Message, client *Client) {
  Todo("Implement user modes")
}
