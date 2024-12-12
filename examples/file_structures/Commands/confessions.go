package commands

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "confessions":
			confessionMessage := ""
			for _, opt := range i.ApplicationCommandData().Options {
				if opt.Name == "message" && opt.Type == discordgo.ApplicationCommandOptionString {
					confessionMessage = opt.StringValue()
					break
				}
			}

			channelID := os.Getenv("ChannelID")
			embed := &discordgo.MessageEmbed{
				Type:        discordgo.EmbedTypeRich,
				Title:       "Confession",
				Description: confessionMessage,
				Color:       0xff0000,
			}

			if confessionMessage != "" {
				_, err := s.ChannelMessageSendEmbed(channelID, embed)
				if err != nil {
					fmt.Println("Error sending embed message: ", err)
				}
			}
		}
	}
}
