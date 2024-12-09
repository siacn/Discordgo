package main

import (
	"flag"
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
			Name:        "poll",
			Description: "Sample command for creating polls",
			Options: []*discordgo.ApplicationCommandOption{
				{
					// Questions are limited to 300 characters
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "question",
					Description: "The question/heading for the poll",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "mulit-select",
					Description: "Allow multiple Choices",
					Required:    true,
				},
				{
					// Duration can go upto 32 days
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "duration",
					Description: "Duration of the poll",
					Required:    true,
					// The duration is an integer indicated in hours
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "1 hour",
							Value: 1,
						},
						{
							Name:  "4 hours",
							Value: 4,
						},
						{
							Name:  "8 hours",
							Value: 8,
						},
						{
							Name:  "24 hours",
							Value: 24,
						},
						{
							Name:  "3 days",
							Value: 72,
						},
						{
							Name:  "1 week",
							Value: 168,
						},
						{
							Name:  "2 weeks",
							Value: 336,
						},
					},
				},
				// Answers and emojis
				{
					// The answers are limited to 55 characters
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer-1",
					Description: "First option for the poll",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer-2",
					Description: "Second option for the poll",
				},
				// Can be further continued upto 10
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"poll": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// Access options in the order provided by the user
			options := i.ApplicationCommandData().Options

			// Create an array of answers
			answers := make([]discordgo.PollAnswer, len(options)-3)

			for i := 0; i < len(answers); i++ {
				answer := discordgo.PollAnswer{
					AnswerID: 0,
					Media: &discordgo.PollMedia{
						Text: options[i+3].StringValue(),
					},
				}
				answers[i] = answer
			}

			// Create a variable to hold the values from the command
			poll := discordgo.Poll{
				Question:         discordgo.PollMedia{Text: i.ApplicationCommandData().Options[0].StringValue()},
				Answers:          answers,
				AllowMultiselect: i.ApplicationCommandData().Options[1].BoolValue(),
				Duration:         int(i.ApplicationCommandData().Options[2].IntValue()),
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Poll: &poll,
				},
			})
		},
	}
)

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

	syncCommands(s, "", commands)

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

}

func syncCommands(s *discordgo.Session, guildID string, desiredCommandList []*discordgo.ApplicationCommand) {
	existingCommands, err := s.ApplicationCommands(s.State.User.ID, guildID)
	if err != nil {
		log.Fatalf("Failed to fetch commands for guild %s: %v", guildID, err)
		return
	}

	desiredMap := make(map[string]*discordgo.ApplicationCommand)
	for _, cmd := range desiredCommandList {
		desiredMap[cmd.Name] = cmd
	}

	existingMap := make(map[string]*discordgo.ApplicationCommand)
	for _, cmd := range existingCommands {
		existingMap[cmd.Name] = cmd
	}

	// Delete commands not in the desired list
	for _, cmd := range existingCommands {
		if _, found := desiredMap[cmd.Name]; !found {
			err := s.ApplicationCommandDelete(s.State.User.ID, guildID, cmd.ID)
			if err != nil {
				log.Printf("Failed to delete command %s (%s) in guild %s: %v", cmd.Name, cmd.ID, guildID, err)
			} else {
				log.Printf("Successfully deleted command %s (%s) in guild %s", cmd.Name, cmd.ID, guildID)
			}
		}
	}

	// Create or update existing commands
	for _, cmd := range desiredCommandList {
		if existingCmd, found := existingMap[cmd.Name]; found {
			// Edit existing command
			_, err := s.ApplicationCommandEdit(s.State.User.ID, guildID, existingCmd.ID, cmd)
			if err != nil {
				log.Printf("Failed to edit command %s (%s) in guild %s: %v", cmd.Name, cmd.ID, guildID, err)
			} else {
				log.Printf("Successfully edited command %s (%s) in guild %s", cmd.Name, cmd.ID, guildID)
			}
		} else {
			// Create new command
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
			if err != nil {
				log.Printf("Failed to create command %s in guild %s: %v", cmd.Name, guildID, err)
			} else {
				log.Printf("Successfully created command %s in guild %s", cmd.Name, guildID)
			}
		}
	}
}
