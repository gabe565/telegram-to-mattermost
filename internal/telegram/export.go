package telegram

import (
	"encoding/json"
	"slices"
)

type Export struct {
	Name     string     `json:"name"`
	Type     ExportType `json:"type"`
	ID       int64      `json:"id"`
	users    []string
	Messages []*Message `json:"messages"`
}

func (e *Export) Users() []string {
	if e.users != nil {
		return slices.Clone(e.users)
	}

	users := make([]string, 0, 2)
	for _, msg := range e.Messages {
		if msg.From != "" && !slices.Contains(users, msg.From) {
			users = append(users, msg.From)
		}
	}
	e.users = slices.Clip(users)
	return e.Users()
}

//go:generate enumer -type ExportType -trimprefix Type -transform snake -text

type ExportType uint8

const (
	TypePersonalChat ExportType = iota
)

type Message struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`

	Date   json.Number  `json:"date_unixtime"`
	Edited *json.Number `json:"edited_unixtime"`

	From             string  `json:"from"`
	FromID           string  `json:"from_id"`
	ForwardedFrom    *string `json:"forwarded_from"`
	ReplyToMessageID *int64  `json:"reply_to_message_id"`
	ViaBot           *string `json:"via_bot"`

	*File
	*Event
	*Location
	*Contact

	TextEntities []*TextEntity `json:"text_entities"`

	Reply    *Message `json:"-"`
	IsPinned *bool

	// Unused fields
	Text         json.RawMessage `json:"text"`
	DateString   json.RawMessage `json:"date"`
	EditedString json.RawMessage `json:"edited"`
}

type File struct {
	Path         *string `json:"file"`
	FileName     *string `json:"file_name"`
	Thumbnail    *string `json:"thumbnail"`
	MediaType    *string `json:"media_type"`
	StickerEmoji *string `json:"sticker_emoji"`
	MIMEType     *string `json:"mime_type"`

	Performer *string `json:"performer"`
	Title     *string `json:"title"`

	Photo *string `json:"photo"`

	DurationSeconds *int `json:"duration_seconds"`

	Width  *int `json:"width"`
	Height *int `json:"height"`
}

type Event struct {
	Actor         *string  `json:"actor"`
	ActorID       *string  `json:"actor_id"`
	Action        *Action  `json:"action"`
	Inviter       *string  `json:"inviter"`
	MessageID     *int64   `json:"message_id"`
	DiscardReason *string  `json:"discard_reason"`
	Members       []string `json:"members"`
}

//go:generate enumer -type Action -trimprefix Action -transform snake -text

type Action uint8

const (
	ActionCreateGroup Action = iota
	ActionEditGroupPhoto
	ActionInviteMembers
	ActionJoinGroupByLink
	ActionPinMessage
	ActionRemoveMembers
	ActionPhoneCall
)

type TextEntity struct {
	Type     TextEntityType `json:"type"`
	Text     string         `json:"text"`
	Href     *string        `json:"href"`
	UserID   *int64         `json:"user_id"`
	Language *string        `json:"language"`
}

type Location struct {
	Location                  *LocationInformation `json:"location_information"`
	LiveLocationPeriodSeconds *int                 `json:"live_location_period_seconds"`
}

type LocationInformation struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type Contact struct {
	ContactInformation *ContactInformation `json:"contact_information"`
	VCard              *string             `json:"contact_vcard"`
}

type ContactInformation struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	PhoneNumber *string `json:"phone_number"`
}
