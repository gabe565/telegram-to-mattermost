package telegram

import (
	"log/slog"
	"strings"
)

//go:generate enumer -type TextEntityType -trimprefix Type -transform snake -text

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
)

func (m *Message) FormatText(maxLen uint) string {
	var n int
	for _, e := range m.TextEntities {
		n += len(e.Text)
		switch e.Type {
		case TypeLink, TypeItalic, TypeCode:
			n += 2
		case TypeTextLink:
			n += 2
			if e.Href != nil {
				n += len(*e.Href)
			}
		case TypeBold:
			n += 4
		case TypeMentionName:
			n += 1
		case TypeEmail:
			n += len(e.Text) + 11
		case TypePhone:
			n += len(e.Text) + 8
		case TypePre:
			if e.Language != nil {
				n += len(*e.Language)
			}
			n += 8
		}
	}

	var buf strings.Builder
	buf.Grow(n)
	for _, e := range m.TextEntities {
		switch e.Type {
		case TypePlain, TypeMention, TypeHashtag:
			buf.WriteString(e.Text)
		case TypeLink:
			buf.WriteByte('<')
			buf.WriteString(e.Text)
			buf.WriteByte('>')
		case TypeTextLink:
			buf.WriteByte('[')
			buf.WriteString(e.Text)
			buf.WriteString("](")
			if e.Href != nil {
				buf.WriteString(*e.Href)
			}
			buf.WriteByte(')')
		case TypeBold:
			buf.WriteString("**")
			buf.WriteString(e.Text)
			buf.WriteString("**")
		case TypeItalic:
			buf.WriteByte('_')
			buf.WriteString(e.Text)
			buf.WriteByte('_')
		case TypeMentionName: // TODO: Handle user ID
			buf.WriteByte('@')
			buf.WriteString(e.Text)
		case TypeEmail:
			buf.WriteByte('[')
			buf.WriteString(e.Text)
			buf.WriteString("](mailto:")
			buf.WriteString(e.Text)
			buf.WriteByte(')')
		case TypePhone:
			buf.WriteByte('[')
			buf.WriteString(e.Text)
			buf.WriteString("](tel:")
			buf.WriteString(e.Text)
			buf.WriteByte(')')
		case TypeCode:
			buf.WriteByte('`')
			buf.WriteString(e.Text)
			buf.WriteByte('`')
		case TypePre:
			buf.WriteString("```")
			if e.Language != nil {
				buf.WriteString(*e.Language)
			}
			buf.WriteByte('\n')
			buf.WriteString(e.Text)
			buf.WriteString("\n```")
		}
	}
	if buf.Len() > int(maxLen) {
		slog.Warn("Truncating message", "length", buf.Len(), "id", m.ID, "from", m.From, "timestamp", m.Date.String())
		return buf.String()[:maxLen]
	}
	return buf.String()
}
