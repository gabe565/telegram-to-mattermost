package telegram

//go:generate enumer -type Action -trimprefix Action -transform snake -text -output action_string.go

type Action uint8

const (
	ActionCreateGroup Action = iota
	ActionEditGroupPhoto
	ActionInviteMembers
	ActionJoinGroupByLink
	ActionPinMessage
	ActionRemoveMembers
	ActionPhoneCall
	ActionScoreInGame
	ActionEditChatTheme
)
