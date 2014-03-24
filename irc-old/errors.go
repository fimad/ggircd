package irc

var (
  ErrorNoSuchNick = Message{
    Command:  "401",
    Trailing: "No such nick",
  }
  ErrorNoSuchServer = Message{
    Command: "402",
    Trailing: "No such server",
  }
  ErrorNoSuchChannel = Message{
    Command: "403",
    Trailing: "No such channel",
  }
  ErrorCannotSendToChan = Message{
    Command: "404",
    Trailing: "Cannot send to channel",
  }
  ErrorTooManyChannels = Message{
    Command: "405",
    Trailing: "Too many channels",
  }
  ErrorWasNoSuchNick = Message{
    Command: "406",
    Trailing: "There was no such nick",
  }
  ErrorTooManyTargets = Message{
    Command: "407",
    Trailing: "Too many targets",
  }
  ErrorNoOrigin = Message{
    Command: "409",
    Trailing: "No origin present",
  }
  ErrorNoRecipient = Message{
    Command: "411",
    Trailing: "No recipient",
  }
  ErrorNoTextToSend = Message{
    Command: "412",
    Trailing: "No text to send",
  }
  ErrorNoTopLevel = Message{
    Command: "413",
    Trailing: "No top level domain",
  }
  ErrorWildTopLevel = Message{
    Command: "414",
    Trailing: "Wild top level domain",
  }
  ErrorUnknownCommand = Message{
    Command: "421",
    Trailing: "Unknown command",
  }
  ErrorNoMOTD = Message{
    Command: "422",
    Trailing: "No MOTD",
  }
  ErrorNoAdminInfo = Message{
    Command: "423",
    Trailing: "No admin info",
  }
  ErrorFileError = Message{
    Command: "424",
    Trailing: "File error",
  }
  ErrorNoNicknameGiven = Message{
    Command: "431",
    Trailing: "No nickname given",
  }
  ErrorErroneusNickname = Message{
    Command: "432",
    Trailing: "Erroneus nickname",
  }
  ErrorNicknameInUse = Message{
    Command: "433",
    Trailing: "Nickname in use",
  }
  ErrorNickCollision = Message{
    Command: "436",
    Trailing: "Nickname collision",
  }
  ErrorUserNotInChannel = Message{
    Command: "441",
    Trailing: "User not in channel",
  }
  ErrorNotOnChannel = Message{
    Command: "442",
    Trailing: "Not on channel",
  }
  ErrorUserOnChannel = Message{
    Command: "443",
    Trailing: "User on channel",
  }
  ErrorNoLogin = Message{
    Command: "444",
    Trailing: "No login",
  }
  ErrorSummonDisabled = Message{
    Command: "445",
    Trailing: "Summon disabled",
  }
  ErrorUsersDisabled = Message{
    Command: "446",
    Trailing: "Users disabled",
  }
  ErrorNotRegistered = Message{
    Command: "451",
    Trailing: "Not registered",
  }
  ErrorNeedMoreParams = Message{
    Command: "461",
    Trailing: "Need more params",
  }
  ErrorAlreadyRegistred = Message{
    Command: "462",
    Trailing: "Already registered",
  }
  ErrorNoPermForHost = Message{
    Command: "463",
    Trailing: "Insufficient permissions for host",
  }
  ErrorPasswdMismatch = Message{
    Command: "464",
    Trailing: "Password mismatch",
  }
  ErrorYoureBannedCreep = Message{
    Command: "465",
    Trailing: "You're banned, creep",
  }
  ErrorKeySet = Message{
    Command: "467",
    Trailing: "Key set",
  }
  ErrorChannelIsFull = Message{
    Command: "471",
    Trailing: "Channel is full",
  }
  ErrorUnknownMode = Message{
    Command: "472",
    Trailing: "Unknown mode",
  }
  ErrorInviteOnlyChan = Message{
    Command: "473",
    Trailing: "Invite only channel",
  }
  ErrorBannedFromChan = Message{
    Command: "474",
    Trailing: "Banned from channel",
  }
  ErrorBadChannelKey = Message{
    Command: "475",
    Trailing: "Bad channel key",
  }
  ErrorNoPrivileges = Message{
    Command: "481",
    Trailing: "No privileges",
  }
  ErrorChanOPrivsNeeded = Message{
    Command: "482",
    Trailing: "Channel +o privileges needed",
  }
  ErrorCantKillServer = Message{
    Command: "483",
    Trailing: "Cannot kill server",
  }
  ErrorNoOperHost = Message{
    Command: "491",
    Trailing: "No operator host",
  }
  ErrorUModeUnknownFlag = Message{
    Command: "501",
    Trailing: "User mode unknown flag",
  }
  ErrorUsersDontMatch = Message{
    Command: "502",
    Trailing: "Users don't match",
  }
)
