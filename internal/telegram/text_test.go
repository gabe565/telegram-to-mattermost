package telegram

import (
	"testing"

	"github.com/gabe565/telegram-to-mattermost/internal/config"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func TestMessage_FormatText(t *testing.T) {
	conf := config.New()

	type fields struct {
		TextEntities []*TextEntity
	}
	type args struct {
		conf *config.Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{"empty", fields{[]*TextEntity{}}, args{conf}, []string{""}},
		{"plain", fields{[]*TextEntity{{Type: TypePlain, Text: "hello"}}}, args{conf}, []string{"hello"}},
		{"multiple", fields{[]*TextEntity{{Type: TypePlain, Text: "hello "}, {Type: TypePlain, Text: "world"}}}, args{conf}, []string{"hello world"}},
		{"mention", fields{[]*TextEntity{{Type: TypeMention, Text: "@gabe565"}}}, args{conf}, []string{"@gabe565"}},
		{"hashtag", fields{[]*TextEntity{{Type: TypeHashtag, Text: "#telegram"}}}, args{conf}, []string{"#telegram"}},
		{"link", fields{[]*TextEntity{{Type: TypeLink, Text: "https://gabecook.com"}}}, args{conf}, []string{"<https://gabecook.com>"}},
		{"text link", fields{[]*TextEntity{{Type: TypeTextLink, Text: "Gabe Cook", Href: ptr.To("https://gabecook.com")}}}, args{conf}, []string{"[Gabe Cook](https://gabecook.com)"}},
		{"text link null href", fields{[]*TextEntity{{Type: TypeTextLink, Text: "Gabe Cook"}}}, args{conf}, []string{"[Gabe Cook]()"}},
		{"strikethrough", fields{[]*TextEntity{{Type: TypeStrikethrough, Text: "Strike"}}}, args{conf}, []string{"~~Strike~~"}},
		{"spoiler", fields{[]*TextEntity{{Type: TypeStrikethrough, Text: "I am your father"}}}, args{conf}, []string{"~~I am your father~~"}},
		{"bold", fields{[]*TextEntity{{Type: TypeBold, Text: "Bold"}}}, args{conf}, []string{"**Bold**"}},
		{"italic", fields{[]*TextEntity{{Type: TypeItalic, Text: "Italic"}}}, args{conf}, []string{"_Italic_"}},
		{"mention name", fields{[]*TextEntity{{Type: TypeMentionName, Text: "Gabe", UserID: ptr.To(int64(1))}}}, args{&config.Config{MaxTextLength: 4000, Users: config.UserList{"user1": {Username: "gabe565"}}}}, []string{"@gabe565"}},
		{"mention name no user match", fields{[]*TextEntity{{Type: TypeMentionName, Text: "Gabe", UserID: ptr.To(int64(1))}}}, args{conf}, []string{"@Gabe"}},
		{"mention name no user ID", fields{[]*TextEntity{{Type: TypeMentionName, Text: "Gabe"}}}, args{conf}, []string{"@Gabe"}},
		{"email", fields{[]*TextEntity{{Type: TypeEmail, Text: "gabe@example.com"}}}, args{conf}, []string{"[gabe@example.com](mailto:gabe@example.com)"}},
		{"phone", fields{[]*TextEntity{{Type: TypePhone, Text: "5178675309"}}}, args{conf}, []string{"[5178675309](tel:5178675309)"}},
		{"code", fields{[]*TextEntity{{Type: TypeCode, Text: "echo hello"}}}, args{conf}, []string{"`echo hello`"}},
		{"pre no language", fields{[]*TextEntity{{Type: TypePre, Text: "echo hello"}}}, args{conf}, []string{"```\necho hello\n```"}},
		{"pre with language", fields{[]*TextEntity{{Type: TypePre, Text: "echo hello", Language: ptr.To("shell")}}}, args{conf}, []string{"```shell\necho hello\n```"}},
		{"split", fields{[]*TextEntity{{Type: TypePlain, Text: "abcdefghijkl"}}}, args{&config.Config{MaxTextLength: 3}}, []string{"abc", "def", "ghi", "jkl"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				DateUnix:     "0",
				TextEntities: tt.fields.TextEntities,
			}
			assert.Equal(t, tt.want, m.FormatText(tt.args.conf))
		})
	}
}
