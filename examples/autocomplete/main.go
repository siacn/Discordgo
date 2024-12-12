package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "single-autocomplete",
			Description: "Showcase of single autocomplete option",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "autocomplete-option",
					Description:  "Autocomplete option",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "multi-autocomplete",
			Description: "Showcase of multiple autocomplete option",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "autocomplete-option-1",
					Description:  "Autocomplete option 1",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:         "autocomplete-option-2",
					Description:  "Autocomplete option 2",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"single-autocomplete": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				data := i.ApplicationCommandData()
				_, err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(
							"You picked %q autocompletion",
							// Autocompleted options do not affect usual flow of handling application command. They are ordinary options at this stage
							data.Options[0].StringValue(),
						),
					},
				})
				if err != nil {
					panic(err)
				}
			// Autocomplete options introduce a new interaction type (8) for returning custom autocomplete results.
			case discordgo.InteractionApplicationCommandAutocomplete:
				data := i.ApplicationCommandData()
				choices := []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Autocomplete",
						Value: "autocomplete",
					},
					{
						Name:  "Autocomplete is best!",
						Value: "autocomplete_is_best",
					},
					{
						Name:  "Choice 3",
						Value: "choice3",
					},
					{
						Name:  "Choice 4",
						Value: "choice4",
					},
					{
						Name:  "Choice 5",
						Value: "choice5",
					},
					// And so on, up to 25 choices
				}

				if data.Options[0].StringValue() != "" {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  data.Options[0].StringValue(), // To get user input you just get value of the autocomplete option.
						Value: "choice_custom",
					})
				}

				_, err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionApplicationCommandAutocompleteResult,
					Data: &discordgo.InteractionResponseData{
						Choices: choices, // This is basically the whole purpose of autocomplete interaction - return custom options to the user.
					},
				})
				if err != nil {
					panic(err)
				}
			}
		},
		"multi-autocomplete": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				data := i.ApplicationCommandData()
				_, err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(
							"Option 1: %s\nOption 2: %s",
							data.Options[0].StringValue(),
							data.Options[1].StringValue(),
						),
					},
				})
				if err != nil {
					panic(err)
				}
			case discordgo.InteractionApplicationCommandAutocomplete:
				data := i.ApplicationCommandData()
				var choices []*discordgo.ApplicationCommandOptionChoice
				switch {
				// In this case there are multiple autocomplete options. The Focused field shows which option user is focused on.
				case data.Options[0].Focused:
					choices = []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Autocomplete 4 first option",
							Value: "autocomplete_default",
						},
						{
							Name:  "Choice 3",
							Value: "choice3",
						},
						{
							Name:  "Choice 4",
							Value: "choice4",
						},
						{
							Name:  "Choice 5",
							Value: "choice5",
						},
					}
					if data.Options[0].StringValue() != "" {
						choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
							Name:  data.Options[0].StringValue(),
							Value: "choice_custom",
						})
					}

				case data.Options[1].Focused:
					choices = []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Autocomplete 4 second option",
							Value: "autocomplete_1_default",
						},
						{
							Name:  "Choice 3.1",
							Value: "choice3_1",
						},
						{
							Name:  "Choice 4.1",
							Value: "choice4_1",
						},
						{
							Name:  "Choice 5.1",
							Value: "choice5_1",
						},
					}
					if data.Options[1].StringValue() != "" {
						choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
							Name:  data.Options[1].StringValue(),
							Value: "choice_custom_2",
						})
					}
				}

				_, err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionApplicationCommandAutocompleteResult,
					Data: &discordgo.InteractionResponseData{
						Choices: choices,
					},
				})
				if err != nil {
					panic(err)
				}
			}
		},
	}
)

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up!") })
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	createdCommands, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, *GuildID, commands)

	if err != nil {
		log.Fatalf("Cannot register commands: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")

	if *RemoveCommands {
		for _, cmd := range createdCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, cmd.ID)
			if err != nil {
				log.Fatalf("Cannot delete %q command: %v", cmd.Name, err)
			}
		}
	}
}
