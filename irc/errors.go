package irc

var (
	errorNoSuchNick = message{
		command:  "401",
		trailing: "No such nick",
	}
	errorNoSuchServer = message{
		command:  "402",
		trailing: "No such server",
	}
	errorNoSuchChannel = message{
		command:  "403",
		trailing: "No such channel",
	}
	errorCannotSendToChan = message{
		command:  "404",
		trailing: "Cannot send to channel",
	}
	errorTooManyChannels = message{
		command:  "405",
		trailing: "Too many channels",
	}
	errorWasNoSuchNick = message{
		command:  "406",
		trailing: "There was no such nick",
	}
	errorTooManyTargets = message{
		command:  "407",
		trailing: "Too many targets",
	}
	errorNoOrigin = message{
		command:  "409",
		trailing: "No origin present",
	}
	errorNoRecipient = message{
		command:  "411",
		trailing: "No recipient",
	}
	errorNoTextToSend = message{
		command:  "412",
		trailing: "No text to send",
	}
	errorNoTopLevel = message{
		command:  "413",
		trailing: "No top level domain",
	}
	errorWildTopLevel = message{
		command:  "414",
		trailing: "Wild top level domain",
	}
	errorUnknowncommand = message{
		command:  "421",
		trailing: "Unknown command",
	}
	errorNoMOTD = message{
		command:  "422",
		trailing: "No MOTD",
	}
	errorNoAdminInfo = message{
		command:  "423",
		trailing: "No admin info",
	}
	errorFileError = message{
		command:  "424",
		trailing: "File error",
	}
	errorNoNicknameGiven = message{
		command:  "431",
		trailing: "No nickname given",
	}
	errorErroneusNickname = message{
		command:  "432",
		trailing: "Erroneus nickname",
	}
	errorNicknameInUse = message{
		command:  "433",
		trailing: "Nickname in use",
	}
	errorNickCollision = message{
		command:  "436",
		trailing: "Nickname collision",
	}
	errorUserNotInChannel = message{
		command:  "441",
		trailing: "User not in channel",
	}
	errorNotOnChannel = message{
		command:  "442",
		trailing: "Not on channel",
	}
	errorUserOnChannel = message{
		command:  "443",
		trailing: "User on channel",
	}
	errorNoLogin = message{
		command:  "444",
		trailing: "No login",
	}
	errorSummonDisabled = message{
		command:  "445",
		trailing: "Summon disabled",
	}
	errorUsersDisabled = message{
		command:  "446",
		trailing: "Users disabled",
	}
	errorNotRegistered = message{
		command:  "451",
		trailing: "Not registered",
	}
	errorNeedMoreParams = message{
		command:  "461",
		trailing: "Need more params",
	}
	errorAlreadyRegistred = message{
		command:  "462",
		trailing: "Already registered",
	}
	errorNoPermForHost = message{
		command:  "463",
		trailing: "Insufficient permissions for host",
	}
	errorPasswdMismatch = message{
		command:  "464",
		trailing: "Password mismatch",
	}
	errorYoureBannedCreep = message{
		command:  "465",
		trailing: "You're banned, creep",
	}
	errorKeySet = message{
		command:  "467",
		trailing: "Key set",
	}
	errorChannelIsFull = message{
		command:  "471",
		trailing: "Channel is full",
	}
	errorUnknownMode = message{
		command:  "472",
		trailing: "Unknown mode",
	}
	errorInviteOnlyChan = message{
		command:  "473",
		trailing: "Invite only channel",
	}
	errorBannedFromChan = message{
		command:  "474",
		trailing: "Banned from channel",
	}
	errorBadChannelKey = message{
		command:  "475",
		trailing: "Bad channel key",
	}
	errorNoPrivileges = message{
		command:  "481",
		trailing: "No privileges",
	}
	errorChanOPrivsNeeded = message{
		command:  "482",
		trailing: "Channel +o privileges needed",
	}
	errorCantKillServer = message{
		command:  "483",
		trailing: "Cannot kill server",
	}
	errorNoOperHost = message{
		command:  "491",
		trailing: "No operator host",
	}
	errorUModeUnknownFlag = message{
		command:  "501",
		trailing: "User mode unknown flag",
	}
	errorUsersDontMatch = message{
		command:  "502",
		trailing: "Users don't match",
	}
)
