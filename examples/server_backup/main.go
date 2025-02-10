package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	BackupDir      = flag.String("backupdir", "backups", "Directory to store backups")
)

var s *discordgo.Session

func init() { flag.Parse() }

// Helper function to convert string to *string
func stringPtr(s string) *string {
	return &s
}

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

// BackupData represents the structure of our backup JSON
type BackupData struct {
	GuildID    string           `json:"guild_id"`
	GuildName  string           `json:"guild_name"`
	BackupTime time.Time        `json:"backup_time"`
	Channels   []*ChannelBackup `json:"channels"`
}

type ChannelBackup struct {
	ID       string                `json:"id"`
	Name     string                `json:"name"`
	Type     discordgo.ChannelType `json:"type"`
	Messages []*discordgo.Message  `json:"messages"`
}

var (
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "backup-server",
			Description:              "Backup all text channels in the server",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
		},
		{
			Name:                     "backup-channel",
			Description:              "Backup a specific channel",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The channel to backup",
					Required:    true,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
					},
				},
			},
		},
		{
			Name:                     "list-backups",
			Description:              "List all available backups",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
		},
		{
			Name:                     "get-backup",
			Description:              "Get a specific backup file",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "filename",
					Description: "The backup file to retrieve",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"backup-server": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Acknowledge the command immediately
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Starting server backup...",
				},
			})

			guild, err := s.Guild(i.GuildID)
			if err != nil {
				sendFollowUpError(s, i, "Failed to get guild information", err)
				return
			}

			backup := &BackupData{
				GuildID:    guild.ID,
				GuildName:  guild.Name,
				BackupTime: time.Now(),
				Channels:   make([]*ChannelBackup, 0),
			}

			channels, err := s.GuildChannels(i.GuildID)
			if err != nil {
				sendFollowUpError(s, i, "Failed to get guild channels", err)
				return
			}

			for _, channel := range channels {
				if channel.Type != discordgo.ChannelTypeGuildText {
					continue
				}

				messages, err := FetchChannelHistory(s, channel.ID)
				if err != nil {
					log.Printf("Error backing up channel %s: %v", channel.Name, err)
					continue
				}

				backup.Channels = append(backup.Channels, &ChannelBackup{
					ID:       channel.ID,
					Name:     channel.Name,
					Type:     channel.Type,
					Messages: messages,
				})
			}

			// Create backup directory if it doesn't exist
			if err := os.MkdirAll(*BackupDir, 0755); err != nil {
				sendFollowUpError(s, i, "Failed to create backup directory", err)
				return
			}

			// Save backup to file
			filename := fmt.Sprintf("%s/%s_%s.json", *BackupDir, guild.Name, time.Now().Format("2006-01-02_15-04-05"))
			data, err := json.MarshalIndent(backup, "", "  ")
			if err != nil {
				sendFollowUpError(s, i, "Failed to marshal backup data", err)
				return
			}

			if err := os.WriteFile(filename, data, 0644); err != nil {
				sendFollowUpError(s, i, "Failed to write backup file", err)
				return
			}

			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: stringPtr(fmt.Sprintf("✅ Server backup completed! Backed up %d channels.\nBackup saved as: `%s`",
					len(backup.Channels), filepath.Base(filename))),
			})
		},
		"backup-channel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			channel := options[0].ChannelValue(s)

			// Acknowledge the command immediately
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Starting backup of channel #%s...", channel.Name),
				},
			})

			messages, err := FetchChannelHistory(s, channel.ID)
			if err != nil {
				sendFollowUpError(s, i, "Failed to fetch channel history", err)
				return
			}

			backup := &BackupData{
				GuildID:    i.GuildID,
				GuildName:  channel.Name,
				BackupTime: time.Now(),
				Channels: []*ChannelBackup{
					{
						ID:       channel.ID,
						Name:     channel.Name,
						Type:     channel.Type,
						Messages: messages,
					},
				},
			}

			// Create backup directory if it doesn't exist
			if err := os.MkdirAll(*BackupDir, 0755); err != nil {
				sendFollowUpError(s, i, "Failed to create backup directory", err)
				return
			}

			// Save backup to file
			filename := fmt.Sprintf("%s/%s_%s.json", *BackupDir, channel.Name, time.Now().Format("2006-01-02_15-04-05"))
			data, err := json.MarshalIndent(backup, "", "  ")
			if err != nil {
				sendFollowUpError(s, i, "Failed to marshal backup data", err)
				return
			}

			if err := os.WriteFile(filename, data, 0644); err != nil {
				sendFollowUpError(s, i, "Failed to write backup file", err)
				return
			}

			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: stringPtr(fmt.Sprintf("✅ Channel backup completed! Backed up %d messages.\nBackup saved as: `%s`",
					len(messages), filepath.Base(filename))),
			})
		},
		"list-backups": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			files, err := os.ReadDir(*BackupDir)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "No backups found or backup directory does not exist.",
					},
				})
				return
			}

			var backupList strings.Builder
			backupList.WriteString("Available backups:\n```\n")
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".json") {
					backupList.WriteString(file.Name() + "\n")
				}
			}
			backupList.WriteString("```")

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: backupList.String(),
				},
			})
		},
		"get-backup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			filename := options[0].StringValue()

			filepath := fmt.Sprintf("%s/%s", *BackupDir, filename)
			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Backup file not found.",
					},
				})
				return
			}

			file, err := os.Open(filepath)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to read backup file.",
					},
				})
				return
			}
			defer file.Close()

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Here's your backup file:",
					Files: []*discordgo.File{
						{
							Name:   filename,
							Reader: file,
						},
					},
				},
			})
		},
	}
)

// FetchChannelHistory retrieves all historical messages from a specific channel.
func FetchChannelHistory(s *discordgo.Session, channelID string) ([]*discordgo.Message, error) {
	var messages []*discordgo.Message
	var lastMessageID string

	for {
		msgs, err := s.ChannelMessages(channelID, 100, lastMessageID, "", "")
		if err != nil {
			return messages, err
		}

		if len(msgs) == 0 {
			break
		}

		messages = append(messages, msgs...)
		lastMessageID = msgs[len(msgs)-1].ID

		if len(msgs) < 100 {
			break
		}

		time.Sleep(300 * time.Millisecond)
	}

	return messages, nil
}

func sendFollowUpError(s *discordgo.Session, i *discordgo.InteractionCreate, message string, err error) {
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: stringPtr(fmt.Sprintf("❌ Error: %s - %v", message, err)),
	})
}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer func() {
		if *RemoveCommands {
			log.Println("Removing commands...")
			for _, v := range registeredCommands {
				err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
				if err != nil {
					log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
				}
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")
}
