package discordgo

import "strings"

// A User stores all data for an individual Discord user.
type User struct {
	// The ID of the user.
	ID string `json:"id"`

	// The email of the user. This is only present when
	// the application possesses the email scope for the user.
	Email string `json:"email"`

	// The user's username.
	Username string `json:"username"`

	// The hash of the user's avatar. Use Session.UserAvatar
	// to retrieve the avatar itself.
	Avatar string `json:"avatar"`

	// The user's chosen language option.
	Locale string `json:"locale"`

	// The discriminator of the user (4 numbers after name).
	Discriminator string `json:"discriminator"`

	// The token of the user. This is only present for
	// the user represented by the current session.
	Token string `json:"token"`

	// Whether the user's email is verified.
	Verified bool `json:"verified"`

	// Whether the user has multi-factor authentication enabled.
	MFAEnabled bool `json:"mfa_enabled"`

	// dm channel with the user, call CreateDM if it doesn't exist
	DMChannel *Channel `json:"dm_channel,omitempty"`

	// Whether the user is a bot.
	Bot bool `json:"bot"`
}

// String returns a unique identifier of the form username#discriminator
func (u *User) String() string {
	return u.Username + "#" + u.Discriminator
}

// Mention return a string which mentions the user
func (u *User) Mention() string {
	return "<@" + u.ID + ">"
}

// AvatarURL returns a URL to the user's avatar.
//    size:    The size of the user's avatar as a power of two
//             if size is an empty string, no size parameter will
//             be added to the URL.
func (u *User) AvatarURL(size string) string {
	var URL string
	if u.Avatar == "" {
		URL = EndpointDefaultUserAvatar(u.Discriminator)
	} else if strings.HasPrefix(u.Avatar, "a_") {
		URL = EndpointUserAvatarAnimated(u.ID, u.Avatar)
	} else {
		URL = EndpointUserAvatar(u.ID, u.Avatar)
	}

	if size != "" {
		return URL + "?size=" + size
	}
	return URL
}

// IsAvatarAnimated indicates if the user has an animated avatar
func (u *User) IsAvatarAnimated() bool {
	return strings.HasPrefix(u.Avatar, "a_")
}

// IsMentionedIn checks if the user is mentioned in the given message
// message      : message to check for mentions
func (u *User) IsMentionedIn(message *Message) bool {
	if message.MentionEveryone {
		return true
	}

	for _, user := range message.Mentions {
		if user.ID == u.ID {
			return true
		}
	}

	return false
}

// CreateDM creates a DM channel between the client and the user if  it is nil,
// populating User.DMChannel with it. Called automagically if DMChannel nil
// when calling SendMessage or SendMessageComplex
func (u *User) CreateDM(s *Session) (err error) {
	if u.DMChannel != nil {
		return
	}

	channel, err := s.UserChannelCreate(u.ID)
	if err == nil {
		u.DMChannel = channel
	}
	return
}

// SendMessage sends a message to the user
// content: message content to send if provided
// embed: embed to attach to the message if provided
// files: files to attach to the message if provided
func (u *User) SendMessage(s *Session, content string, embed *MessageEmbed, files []*File) (message *Message, err error) {
	if u.DMChannel == nil {
		err = u.CreateDM(s)
		if err != nil {
			return
		}
	}

	return u.DMChannel.SendMessage(s, content, embed, files)
}

// SendMessageComplex sends a message to the user
// data: MessageSend object with the data to send
func (u *User) SendMessageComplex(s *Session, data *MessageSend) (message *Message, err error) {
	if u.DMChannel == nil {
		err = u.CreateDM(s)
		if err != nil {
			return
		}
	}

	return u.DMChannel.SendMessageComplex(s, data)
}

// GetHistory fetches up to limit messages from the user
// limit     : The number messages that can be returned. (max 100)
// beforeID  : If provided all messages returned will be before given ID.
// afterID   : If provided all messages returned will be after given ID.
// aroundID  : If provided all messages returned will be around given ID.
func (u *User) GetHistory(s *Session, limit int, beforeID, afterID, aroundID string) (st []*Message, err error) {
	return s.ChannelMessages(u.DMChannel.ID, limit, beforeID, afterID, aroundID)
}
