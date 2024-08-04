package telegram

//go:generate enumer -type TextEntityType -trimprefix Type -transform snake -text -output text_entity_type_string.go

type TextEntityType uint8

const (
	TypePlain TextEntityType = iota
	TypeLink
	TypeTextLink
	TypeBold
	TypeHashtag
	TypeItalic
	TypeMention
	TypeMentionName
	TypeEmail
	TypePhone
	TypeCode
	TypePre
	TypeStrikethrough
)
