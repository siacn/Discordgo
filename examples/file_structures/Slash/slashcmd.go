package Slash

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func RegisterSlashCommand(s *discordgo.Session, applicationID string) error {
	command := &discordgo.ApplicationCommand{
		Name:        "confessions",
		Description: "Request confessions from the bot with an optional message.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "Your confession message",
				Required:    false,
			},
		},
	}

	_, err := s.ApplicationCommandCreate(applicationID, "", command)
	if err != nil {
		return fmt.Errorf("error registering slash command: %v", err)
	}

	fmt.Println("Slash command '/confessions' registered successfully.")
	return nil
}
