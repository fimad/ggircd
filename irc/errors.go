package irc

var (
  ErrorNoSuchNick       = Message{Command: "401"}
  ErrorNoSuchServer     = Message{Command: "402"}
  ErrorNoSuchChannel    = Message{Command: "403"}
  ErrorCannotSendToChan = Message{Command: "404"}
  ErrorTooManyChannels  = Message{Command: "405"}
  ErrorWasNoSuchNick    = Message{Command: "406"}
  ErrorTooManyTargets   = Message{Command: "407"}
  ErrorNoOrigin         = Message{Command: "409"}
  ErrorNoRecipient      = Message{Command: "411"}
  ErrorNoTextToSend     = Message{Command: "412"}
  ErrorNoTopLevel       = Message{Command: "413"}
  ErrorWildTopLevel     = Message{Command: "414"}
  ErrorUnknownCommand   = Message{Command: "421"}
  ErrorNoMOTD           = Message{Command: "422"}
  ErrorNoAdminInfo      = Message{Command: "423"}
  ErrorFileError        = Message{Command: "424"}
  ErrorNoNicknameGiven  = Message{Command: "431"}
  ErrorErroneusNickname = Message{Command: "432"}
  ErrorNicknameInUse    = Message{Command: "433"}
  ErrorNickCollision    = Message{Command: "436"}
  ErrorUserNotInChannel = Message{Command: "441"}
  ErrorNotOnChannel     = Message{Command: "442"}
  ErrorUserOnChannel    = Message{Command: "443"}
  ErrorNoLogin          = Message{Command: "444"}
  ErrorSummonDisabled   = Message{Command: "445"}
  ErrorUsersDisabled    = Message{Command: "446"}
  ErrorNotRegistered    = Message{Command: "451"}
  ErrorNeedMoreParams   = Message{Command: "461"}
  ErrorAlreadyRegistred = Message{Command: "462"}
  ErrorNoPermForHost    = Message{Command: "463"}
  ErrorPasswdMismatch   = Message{Command: "464"}
  ErrorYoureBannedCreep = Message{Command: "465"}
  ErrorKeySet           = Message{Command: "467"}
  ErrorChannelIsFull    = Message{Command: "471"}
  ErrorUnknownMode      = Message{Command: "472"}
  ErrorInviteOnlyChan   = Message{Command: "473"}
  ErrorBannedFromChan   = Message{Command: "474"}
  ErrorBadChannelKEY    = Message{Command: "475"}
  ErrorNoPrivileges     = Message{Command: "481"}
  ErrorChanOPrivsNeeded = Message{Command: "482"}
  ErrorCantKillServer   = Message{Command: "483"}
  ErrorNoOperHost       = Message{Command: "491"}
  ErrorUModeUnknownFlag = Message{Command: "501"}
  ErrorUsersDontMatch   = Message{Command: "502"}
)
