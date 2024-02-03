// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains variables for all known Discord end points.  All functions
// throughout the Discordgo package use these variables for all connections
// to Discord.  These are all exported and you may modify them if needed.

package discordgo

import (
	"strconv"
)

// APIVersion is the Discord API version used for the REST and Websocket API.
var APIVersion = "9"

// Known Discord API Endpoints.
var (
	EndpointStatus     = "https://status.discord.com/api/v2/"
	EndpointSm         = EndpointStatus + "scheduled-maintenances/"
	EndpointSmActive   = EndpointSm + "active.json"
	EndpointSmUpcoming = EndpointSm + "upcoming.json"

	EndpointDiscord        = "https://discord.com/"
	EndpointAPI            = EndpointDiscord + "api/v" + APIVersion + "/"
	EndpointGuilds         = EndpointAPI + "guilds/"
	EndpointChannels       = EndpointAPI + "channels/"
	EndpointUsers          = EndpointAPI + "users/"
	EndpointGateway        = EndpointAPI + "gateway"
	EndpointGatewayBot     = EndpointGateway + "/bot"
	EndpointWebhooks       = EndpointAPI + "webhooks/"
	EndpointStickers       = EndpointAPI + "stickers/"
	EndpointStageInstances = EndpointAPI + "stage-instances"

	EndpointCDN             = "https://cdn.discordapp.com/"
	EndpointCDNAttachments  = EndpointCDN + "attachments/"
	EndpointCDNAvatars      = EndpointCDN + "avatars/"
	EndpointCDNIcons        = EndpointCDN + "icons/"
	EndpointCDNSplashes     = EndpointCDN + "splashes/"
	EndpointCDNChannelIcons = EndpointCDN + "channel-icons/"
	EndpointCDNBanners      = EndpointCDN + "banners/"
	EndpointCDNGuilds       = EndpointCDN + "guilds/"
	EndpointCDNRoleIcons    = EndpointCDN + "role-icons/"

	EndpointVoice        = EndpointAPI + "/voice/"
	EndpointVoiceRegions = EndpointVoice + "regions"

	EndpointUser       = func(uID Snowflake) string { return EndpointUsers + string(uID) }
	EndpointUserAvatar = func(uID Snowflake, hash string) string {
		return EndpointCDNAvatars + string(uID) + "/" + hash + ".png"
	}
	EndpointUserAvatarAnimated = func(uID Snowflake, hash string) string {
		return EndpointCDNAvatars + string(uID) + "/" + hash + ".gif"
	}
	EndpointDefaultUserAvatar = func(idx int) string {
		return EndpointCDN + "embed/avatars/" + strconv.Itoa(idx) + ".png"
	}
	EndpointUserBanner = func(uID Snowflake, hash string) string {
		return EndpointCDNBanners + string(uID) + "/" + hash + ".png"
	}
	EndpointUserBannerAnimated = func(uID Snowflake, hash string) string {
		return EndpointCDNBanners + string(uID) + "/" + hash + ".gif"
	}

	EndpointUserGuilds                    = func(uID Snowflake) string { return EndpointUsers + string(uID) + "/guilds" }
	EndpointUserGuild                     = func(uID, gID Snowflake) string { return EndpointUsers + string(uID) + "/guilds/" + string(gID) }
	EndpointUserGuildMember               = func(uID, gID Snowflake) string { return EndpointUserGuild(uID, gID) + "/member" }
	EndpointUserChannels                  = func(uID Snowflake) string { return EndpointUsers + string(uID) + "/channels" }
	EndpointUserApplicationRoleConnection = func(aID Snowflake) string {
		return EndpointUsers + "@me/applications/" + string(aID) + "/role-connection"
	}
	EndpointUserConnections = func(uID Snowflake) string { return EndpointUsers + string(uID) + "/connections" }

	EndpointGuild                    = func(gID Snowflake) string { return EndpointGuilds + string(gID) }
	EndpointGuildAutoModeration      = func(gID Snowflake) string { return EndpointGuild(gID) + "/auto-moderation" }
	EndpointGuildAutoModerationRules = func(gID Snowflake) string { return EndpointGuildAutoModeration(gID) + "/rules" }
	EndpointGuildAutoModerationRule  = func(gID, rID Snowflake) string { return EndpointGuildAutoModerationRules(gID) + "/" + string(rID) }
	EndpointGuildThreads             = func(gID Snowflake) string { return EndpointGuild(gID) + "/threads" }
	EndpointGuildActiveThreads       = func(gID Snowflake) string { return EndpointGuildThreads(gID) + "/active" }
	EndpointGuildPreview             = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/preview" }
	EndpointGuildChannels            = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/channels" }
	EndpointGuildMembers             = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/members" }
	EndpointGuildMembersSearch       = func(gID Snowflake) string { return EndpointGuildMembers(gID) + "/search" }
	EndpointGuildMember              = func(gID, uID Snowflake) string { return EndpointGuilds + string(gID) + "/members/" + string(uID) }
	EndpointGuildMemberRole          = func(gID, uID, rID Snowflake) string {
		return EndpointGuilds + string(gID) + "/members/" + string(uID) + "/roles/" + string(rID)
	}
	EndpointGuildBans         = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/bans" }
	EndpointGuildBan          = func(gID, uID Snowflake) string { return EndpointGuilds + string(gID) + "/bans/" + string(uID) }
	EndpointGuildIntegrations = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/integrations" }
	EndpointGuildIntegration  = func(gID, iID Snowflake) string {
		return EndpointGuilds + string(gID) + "/integrations/" + string(iID)
	}
	EndpointGuildRoles        = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/roles" }
	EndpointGuildRole         = func(gID, rID Snowflake) string { return EndpointGuilds + string(gID) + "/roles/" + string(rID) }
	EndpointGuildInvites      = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/invites" }
	EndpointGuildWidget       = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/widget" }
	EndpointGuildEmbed        = EndpointGuildWidget
	EndpointGuildPrune        = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/prune" }
	EndpointGuildIcon         = func(gID Snowflake, hash string) string { return EndpointCDNIcons + string(gID) + "/" + hash + ".png" }
	EndpointGuildIconAnimated = func(gID Snowflake, hash string) string { return EndpointCDNIcons + string(gID) + "/" + hash + ".gif" }
	EndpointGuildSplash       = func(gID Snowflake, hash string) string {
		return EndpointCDNSplashes + string(gID) + "/" + hash + ".png"
	}
	EndpointGuildWebhooks  = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/webhooks" }
	EndpointGuildAuditLogs = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/audit-logs" }
	EndpointGuildEmojis    = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/emojis" }
	EndpointGuildEmoji     = func(gID, eID Snowflake) string { return EndpointGuilds + string(gID) + "/emojis/" + string(eID) }
	EndpointGuildBanner    = func(gID Snowflake, hash string) string {
		return EndpointCDNBanners + string(gID) + "/" + hash + ".png"
	}
	EndpointGuildBannerAnimated = func(gID Snowflake, hash string) string {
		return EndpointCDNBanners + string(gID) + "/" + hash + ".gif"
	}
	EndpointGuildStickers        = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/stickers" }
	EndpointGuildSticker         = func(gID, sID Snowflake) string { return EndpointGuilds + string(gID) + "/stickers/" + string(sID) }
	EndpointStageInstance        = func(cID Snowflake) string { return EndpointStageInstances + "/" + string(cID) }
	EndpointGuildScheduledEvents = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/scheduled-events" }
	EndpointGuildScheduledEvent  = func(gID, eID Snowflake) string {
		return EndpointGuilds + string(gID) + "/scheduled-events/" + string(eID)
	}
	EndpointGuildScheduledEventUsers = func(gID, eID Snowflake) string { return EndpointGuildScheduledEvent(gID, eID) + "/users" }
	EndpointGuildOnboarding          = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/onboarding" }
	EndpointGuildTemplate            = func(code string) string { return EndpointGuilds + "templates/" + code }
	EndpointGuildTemplates           = func(gID Snowflake) string { return EndpointGuilds + string(gID) + "/templates" }
	EndpointGuildTemplateSync        = func(gID Snowflake, code string) string { return EndpointGuilds + string(gID) + "/templates/" + code }
	EndpointGuildMemberAvatar        = func(gID, uID Snowflake, hash string) string {
		return EndpointCDNGuilds + string(gID) + "/users/" + string(uID) + "/avatars/" + hash + ".png"
	}
	EndpointGuildMemberAvatarAnimated = func(gID, uID Snowflake, hash string) string {
		return EndpointCDNGuilds + string(gID) + "/users/" + string(uID) + "/avatars/" + hash + ".gif"
	}

	EndpointRoleIcon = func(rID, hash string) string {
		return EndpointCDNRoleIcons + string(rID) + "/" + hash + ".png"
	}

	EndpointChannel                             = func(cID Snowflake) string { return EndpointChannels + string(cID) }
	EndpointChannelThreads                      = func(cID Snowflake) string { return EndpointChannel(cID) + "/threads" }
	EndpointChannelActiveThreads                = func(cID Snowflake) string { return EndpointChannelThreads(cID) + "/active" }
	EndpointChannelPublicArchivedThreads        = func(cID Snowflake) string { return EndpointChannelThreads(cID) + "/archived/public" }
	EndpointChannelPrivateArchivedThreads       = func(cID Snowflake) string { return EndpointChannelThreads(cID) + "/archived/private" }
	EndpointChannelJoinedPrivateArchivedThreads = func(cID Snowflake) string { return EndpointChannel(cID) + "/users/@me/threads/archived/private" }
	EndpointChannelPermissions                  = func(cID Snowflake) string { return EndpointChannels + string(cID) + "/permissions" }
	EndpointChannelPermission                   = func(cID, tID Snowflake) string {
		return EndpointChannels + string(cID) + "/permissions/" + string(tID)
	}
	EndpointChannelInvites            = func(cID Snowflake) string { return EndpointChannels + string(cID) + "/invites" }
	EndpointChannelTyping             = func(cID Snowflake) string { return EndpointChannels + string(cID) + "/typing" }
	EndpointChannelMessages           = func(cID Snowflake) string { return EndpointChannels + string(cID) + "/messages" }
	EndpointChannelMessage            = func(cID, mID Snowflake) string { return EndpointChannels + string(cID) + "/messages/" + string(mID) }
	EndpointChannelMessageThread      = func(cID, mID Snowflake) string { return EndpointChannelMessage(cID, mID) + "/threads" }
	EndpointChannelMessagesBulkDelete = func(cID Snowflake) string { return EndpointChannel(cID) + "/messages/bulk-delete" }
	EndpointChannelMessagesPins       = func(cID Snowflake) string { return EndpointChannel(cID) + "/pins" }
	EndpointChannelMessagePin         = func(cID, mID Snowflake) string { return EndpointChannel(cID) + "/pins/" + string(mID) }
	EndpointChannelMessageCrosspost   = func(cID, mID Snowflake) string {
		return EndpointChannel(cID) + "/messages/" + string(mID) + "/crosspost"
	}
	EndpointChannelFollow = func(cID Snowflake) string { return EndpointChannel(cID) + "/followers" }
	EndpointThreadMembers = func(tID Snowflake) string { return EndpointChannel(tID) + "/thread-members" }
	EndpointThreadMember  = func(tID, mID Snowflake) string { return EndpointThreadMembers(tID) + "/" + string(mID) }

	EndpointGroupIcon = func(cID Snowflake, hash string) string {
		return EndpointCDNChannelIcons + string(cID) + "/" + hash + ".png"
	}

	EndpointSticker            = func(sID Snowflake) string { return EndpointStickers + string(sID) }
	EndpointNitroStickersPacks = EndpointAPI + "/sticker-packs"

	EndpointChannelWebhooks = func(cID Snowflake) string { return EndpointChannel(cID) + "/webhooks" }
	EndpointWebhook         = func(wID Snowflake) string { return EndpointWebhooks + string(wID) }
	EndpointWebhookToken    = func(wID Snowflake, token string) string { return EndpointWebhooks + string(wID) + "/" + token }
	EndpointWebhookMessage  = func(wID Snowflake, token string, messageID Snowflake) string {
		return EndpointWebhookToken(wID, token) + "/messages/" + string(messageID)
	}

	EndpointMessageReactionsAll = func(cID, mID Snowflake) string {
		return EndpointChannelMessage(cID, mID) + "/reactions"
	}
	EndpointMessageReactions = func(cID, mID, eID Snowflake) string {
		return EndpointChannelMessage(cID, mID) + "/reactions/" + string(eID)
	}
	EndpointMessageReaction = func(cID, mID, eID, uID Snowflake) string {
		return EndpointMessageReactions(cID, mID, eID) + "/" + string(uID)
	}

	EndpointApplicationGlobalCommands = func(aID Snowflake) string {
		return EndpointApplication(aID) + "/commands"
	}
	EndpointApplicationGlobalCommand = func(aID, cID Snowflake) string {
		return EndpointApplicationGlobalCommands(aID) + "/" + string(cID)
	}

	EndpointApplicationGuildCommands = func(aID, gID Snowflake) string {
		return EndpointApplication(aID) + "/guilds/" + string(gID) + "/commands"
	}
	EndpointApplicationGuildCommand = func(aID, gID, cID Snowflake) string {
		return EndpointApplicationGuildCommands(aID, gID) + "/" + string(cID)
	}
	EndpointApplicationCommandPermissions = func(aID, gID, cID Snowflake) string {
		return EndpointApplicationGuildCommand(aID, gID, cID) + "/permissions"
	}
	EndpointApplicationCommandsGuildPermissions = func(aID, gID Snowflake) string {
		return EndpointApplicationGuildCommands(aID, gID) + "/permissions"
	}
	EndpointInteraction = func(aID Snowflake, iToken string) string {
		return EndpointAPI + "interactions/" + string(aID) + "/" + iToken
	}
	EndpointInteractionResponse = func(iID Snowflake, iToken string) string {
		return EndpointInteraction(iID, iToken) + "/callback"
	}
	EndpointInteractionResponseActions = func(aID Snowflake, iToken string) string {
		return EndpointWebhookMessage(aID, iToken, "@original")
	}
	EndpointFollowupMessage = func(aID Snowflake, iToken string) string {
		return EndpointWebhookToken(aID, iToken)
	}
	EndpointFollowupMessageActions = func(aID Snowflake, iToken string, mID Snowflake) string {
		return EndpointWebhookMessage(aID, iToken, mID)
	}

	EndpointGuildCreate = EndpointAPI + "guilds"

	EndpointInvite = func(code string) string { return EndpointAPI + "invites/" + code }

	EndpointEmoji         = func(eID Snowflake) string { return EndpointCDN + "emojis/" + string(eID) + ".png" }
	EndpointEmojiAnimated = func(eID Snowflake) string { return EndpointCDN + "emojis/" + string(eID) + ".gif" }

	EndpointApplications                      = EndpointAPI + "applications"
	EndpointApplication                       = func(aID Snowflake) string { return EndpointApplications + "/" + string(aID) }
	EndpointApplicationRoleConnectionMetadata = func(aID Snowflake) string { return EndpointApplication(aID) + "/role-connections/metadata" }

	EndpointOAuth2                  = EndpointAPI + "oauth2/"
	EndpointOAuth2Applications      = EndpointOAuth2 + "applications"
	EndpointOAuth2Application       = func(aID Snowflake) string { return EndpointOAuth2Applications + "/" + string(aID) }
	EndpointOAuth2ApplicationsBot   = func(aID Snowflake) string { return EndpointOAuth2Applications + "/" + string(aID) + "/bot" }
	EndpointOAuth2ApplicationAssets = func(aID Snowflake) string { return EndpointOAuth2Applications + "/" + string(aID) + "/assets" }

	// TODO: Deprecated, remove in the next release
	EndpointOauth2                  = EndpointOAuth2
	EndpointOauth2Applications      = EndpointOAuth2Applications
	EndpointOauth2Application       = EndpointOAuth2Application
	EndpointOauth2ApplicationsBot   = EndpointOAuth2ApplicationsBot
	EndpointOauth2ApplicationAssets = EndpointOAuth2ApplicationAssets
)
