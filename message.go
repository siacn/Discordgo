// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains code related to the Message struct

package discordgo

import (
	"io"
	"regexp"
)

// A Message stores all data related to a specific Discord message.
type Message struct {
	ID              string               `json:"id"`
	ChannelID       string               `json:"channel_id"`
	Content         string               `json:"content"`
	Timestamp       Timestamp            `json:"timestamp"`
	EditedTimestamp Timestamp            `json:"edited_timestamp"`
	MentionRoles    []string             `json:"mention_roles"`
	Tts             bool                 `json:"tts"`
	MentionEveryone bool                 `json:"mention_everyone"`
	Author          *User                `json:"author"`
	Attachments     []*MessageAttachment `json:"attachments"`
	Embeds          []*MessageEmbed      `json:"embeds"`
	Mentions        []*User              `json:"mentions"`
	Reactions       []*MessageReactions  `json:"reactions"`
}

// File stores info about files you e.g. send in messages.
type File struct {
	Name   string
	Reader io.Reader
}

// MessageSend stores all parameters you can send with ChannelMessageSendComplex.
type MessageSend struct {
	Content string        `json:"content,omitempty"`
	Embed   *MessageEmbed `json:"embed,omitempty"`
	Tts     bool          `json:"tts"`
	File    *File         `json:"file"`
}

// MessageEdit is used to chain parameters via ChannelMessageEditComplex, which
// is also where you should get the instance from.
type MessageEdit struct {
	Content *string       `json:"content,omitempty"`
	Embed   *MessageEmbed `json:"embed,omitempty"`

	ID      string
	Channel string
}

// NewMessageEdit returns a MessageEdit struct, initialized
// with the Channel and ID.
func NewMessageEdit(channelID string, messageID string) *MessageEdit {
	return &MessageEdit{
		Channel: channelID,
		ID:      messageID,
	}
}

// SetContent is the same as setting the variable Content,
// except it doesn't take a pointer.
func (m *MessageEdit) SetContent(str string) *MessageEdit {
	m.Content = &str
	return m
}

// SetEmbed is a convenience function for setting the embed,
// so you can chain commands.
func (m *MessageEdit) SetEmbed(embed *MessageEmbed) *MessageEdit {
	m.Embed = embed
	return m
}

// A MessageAttachment stores data for message attachments.
type MessageAttachment struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Filename string `json:"filename"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
}

// MessageEmbedFooter is a part of a MessageEmbed struct.
type MessageEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedImage is a part of a MessageEmbed struct.
type MessageEmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedThumbnail is a part of a MessageEmbed struct.
type MessageEmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedVideo is a part of a MessageEmbed struct.
type MessageEmbedVideo struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedProvider is a part of a MessageEmbed struct.
type MessageEmbedProvider struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

// MessageEmbedAuthor is a part of a MessageEmbed struct.
type MessageEmbedAuthor struct {
	URL          string `json:"url,omitempty"`
	Name         string `json:"name,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedField is a part of a MessageEmbed struct.
type MessageEmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// An MessageEmbed stores data for message embeds.
type MessageEmbed struct {
	URL         string                 `json:"url,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Color       int                    `json:"color,omitempty"`
	Footer      *MessageEmbedFooter    `json:"footer,omitempty"`
	Image       *MessageEmbedImage     `json:"image,omitempty"`
	Thumbnail   *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *MessageEmbedVideo     `json:"video,omitempty"`
	Provider    *MessageEmbedProvider  `json:"provider,omitempty"`
	Author      *MessageEmbedAuthor    `json:"author,omitempty"`
	Fields      []*MessageEmbedField   `json:"fields,omitempty"`
}

// MessageReactions holds a reactions object for a message.
type MessageReactions struct {
	Count int    `json:"count"`
	Me    bool   `json:"me"`
	Emoji *Emoji `json:"emoji"`
}

// ContentWithMentionsReplaced will replace all @<id> mentions with the
// username of the mention.
func (m *Message) ContentWithMentionsReplaced() (content string) {
	content = m.Content

	for _, user := range m.Mentions {
		content = regexp.MustCompile("<@!?"+regexp.QuoteMeta(user.ID)+">").ReplaceAllString(content, "@"+user.Username)
	}
	return
}

// ContentWithMoreMentionsReplaced will replace all @<id> mentions with the
// username of the mention, but also role IDs and more.
func (m *Message) ContentWithMoreMentionsReplaced(s *Session) (content string, err error) {
	content = m.Content

	if !s.StateEnabled {
		content = m.ContentWithMentionsReplaced()
		return
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		return
	}

	for _, user := range m.Mentions {
		member, err := s.State.Member(channel.GuildID, user.ID)
		if err != nil {
			continue
		}

		nick := member.Nick
		if nick == "" {
			nick = user.Username
		}
		content = regexp.MustCompile("<@!?"+regexp.QuoteMeta(user.ID)+">").ReplaceAllString(content, "@"+nick)
	}
	for _, role := range guild.Roles {
		if !role.Mentionable {
			continue
		}
		content = regexp.MustCompile("<@&"+regexp.QuoteMeta(role.ID)+">").ReplaceAllString(content, "@"+role.Name)
	}
	return
}
