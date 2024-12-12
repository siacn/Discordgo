package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID  = flag.String("guild", "", "Test guild ID")
	BotToken = flag.String("token", "", "Bot access token")
	AppID    = flag.String("app", "", "Application ID")
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

// Important note: call every command in order it's placed in the example.

var (
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"fd_no": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			_, err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Huh. I see, maybe some of these resources might help you?",
					Flags:   discordgo.MessageFlagsEphemeral,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: &discordgo.ComponentEmoji{
										Name: "📜",
									},
									Label: "Documentation",
									Style: discordgo.LinkButton,
									URL:   "https://discord.com/developers/docs/interactions/message-components#buttons",
								},
								discordgo.Button{
									Emoji: &discordgo.ComponentEmoji{
										Name: "🔧",
									},
									Label: "Discord developers",
									Style: discordgo.LinkButton,
									URL:   "https://discord.gg/discord-developers",
								},
								discordgo.Button{
									Emoji: &discordgo.ComponentEmoji{
										Name: "🦫",
									},
									Label: "Discord Gophers",
									Style: discordgo.LinkButton,
									URL:   "https://discord.gg/7RuRrVHyXF",
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"fd_yes": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			_, err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Great! If you wanna know more or just have questions, feel free to visit Discord Devs and Discord Gophers server. " +
						"But now, when you know how buttons work, let's move onto select menus (execute `/selects single`)",
					Flags: discordgo.MessageFlagsEphemeral,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: &discordgo.ComponentEmoji{
										Name: "🔧",
									},
									Label: "Discord developers",
									Style: discordgo.LinkButton,
									URL:   "https://discord.gg/discord-developers",
								},
								discordgo.Button{
									Emoji: &discordgo.ComponentEmoji{
										Name: "🦫",
									},
									Label: "Discord Gophers",
									Style: discordgo.LinkButton,
									URL:   "https://discord.gg/7RuRrVHyXF",
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"select": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var response *discordgo.InteractionResponse

			data := i.MessageComponentData()
			switch data.Values[0] {
			case "go":
				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "This is the way.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				}
			default:
				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "It is not the way to go.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				}
			}
			_, err := s.InteractionRespond(i.Interaction, response)
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second) // Doing that so user won't see instant response.
			_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Anyways, now when you know how to use single select menus, let's see how multi select menus work. " +
					"Try calling `/selects multi` command.",
				Flags: discordgo.MessageFlagsEphemeral,
			})
			if err != nil {
				panic(err)
			}
		},
		"stackoverflow_tags": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			data := i.MessageComponentData()

			const stackoverflowFormat = `https://stackoverflow.com/questions/tagged/%s`

			_, err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Here is your stackoverflow URL: " + fmt.Sprintf(stackoverflowFormat, strings.Join(data.Values, "+")),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second) // Doing that so user won't see instant response.
			_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "But wait, there is more! You can also auto populate the select menu. Try executing `/selects auto-populated`.",
				Flags:   discordgo.MessageFlagsEphemeral,
			})
			if err != nil {
				panic(err)
			}
		},
		"channel_select": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			_, err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This is it. You've reached your destination. Your choice was <#" + i.MessageComponentData().Values[0] + ">\n" +
						"If you want to know more, check out the links below",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: &discordgo.ComponentEmoji{
										Name: "📜",
									},
									Label: "Documentation",
									Style: discordgo.LinkButton,
									URL:   "https://discord.com/developers/docs/interactions/message-components#select-menus",
								},
								discordgo.Button{
									Emoji: &discordgo.ComponentEmoji{
										Name: "🔧",
									},
									Label: "Discord developers",
									Style: discordgo.LinkButton,
									URL:   "https://discord.gg/discord-developers",
								},
								discordgo.Button{
									Emoji: &discordgo.ComponentEmoji{
										Name: "🦫",
									},
									Label: "Discord Gophers",
									Style: discordgo.LinkButton,
									URL:   "https://discord.gg/7RuRrVHyXF",
								},
							},
						},
					},

					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}
		},
	}
	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"buttons": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			_, err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Are you comfortable with buttons and other message components?",
					Flags:   discordgo.MessageFlagsEphemeral,
					// Buttons and other components are specified in Components field.
					Components: []discordgo.MessageComponent{
						// ActionRow is a container of all buttons within the same row.
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									// Label is what the user will see on the button.
									Label: "Yes",
									// Style provides coloring of the button. There are not so many styles tho.
									Style: discordgo.SuccessButton,
									// Disabled allows bot to disable some buttons for users.
									Disabled: false,
									// CustomID is a thing telling Discord which data to send when this button will be pressed.
									CustomID: "fd_yes",
								},
								discordgo.Button{
									Label:    "No",
									Style:    discordgo.DangerButton,
									Disabled: false,
									CustomID: "fd_no",
								},
								discordgo.Button{
									Label:    "I don't know",
									Style:    discordgo.LinkButton,
									Disabled: false,
									// Link buttons don't require CustomID and do not trigger the gateway/HTTP event
									URL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
									Emoji: &discordgo.ComponentEmoji{
										Name: "🤷",
									},
								},
							},
						},
						// The message may have multiple actions rows.
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Label:    "Discord Developers server",
									Style:    discordgo.LinkButton,
									Disabled: false,
									URL:      "https://discord.gg/discord-developers",
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"selects": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var response *discordgo.InteractionResponse
			switch i.ApplicationCommandData().Options[0].Name {
			case "single":
				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Now let's take a look on selects. This is single item select menu.",
						Flags:   discordgo.MessageFlagsEphemeral,
						Components: []discordgo.MessageComponent{
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.SelectMenu{
										// Select menu, as other components, must have a customID, so we set it to this value.
										CustomID:    "select",
										Placeholder: "Choose your favorite programming language 👇",
										Options: []discordgo.SelectMenuOption{
											{
												Label: "Go",
												// As with components, this things must have their own unique "id" to identify which is which.
												// In this case such id is Value field.
												Value: "go",
												Emoji: &discordgo.ComponentEmoji{
													Name: "🦦",
												},
												// You can also make it a default option, but in this case we won't.
												Default:     false,
												Description: "Go programming language",
											},
											{
												Label: "JS",
												Value: "js",
												Emoji: &discordgo.ComponentEmoji{
													Name: "🟨",
												},
												Description: "JavaScript programming language",
											},
											{
												Label: "Python",
												Value: "py",
												Emoji: &discordgo.ComponentEmoji{
													Name: "🐍",
												},
												Description: "Python programming language",
											},
										},
									},
								},
							},
						},
					},
				}
			case "multi":
				minValues := 1
				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Now let's see how the multi-item select menu works: " +
							"try generating your own stackoverflow search link",
						Flags: discordgo.MessageFlagsEphemeral,
						Components: []discordgo.MessageComponent{
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.SelectMenu{
										CustomID:    "stackoverflow_tags",
										Placeholder: "Select tags to search on StackOverflow",
										// This is where confusion comes from. If you don't specify these things you will get single item select.
										// These fields control the minimum and maximum amount of selected items.
										MinValues: &minValues,
										MaxValues: 3,
										Options: []discordgo.SelectMenuOption{
											{
												Label:       "Go",
												Description: "Simple yet powerful programming language",
												Value:       "go",
												// Default works the same for multi-select menus.
												Default: false,
												Emoji: &discordgo.ComponentEmoji{
													Name: "🦦",
												},
											},
											{
												Label:       "JS",
												Description: "Multiparadigm OOP language",
												Value:       "javascript",
												Emoji: &discordgo.ComponentEmoji{
													Name: "🟨",
												},
											},
											{
												Label:       "Python",
												Description: "OOP prototyping programming language",
												Value:       "python",
												Emoji: &discordgo.ComponentEmoji{
													Name: "🐍",
												},
											},
											{
												Label:       "Web",
												Description: "Web related technologies",
												Value:       "web",
												Emoji: &discordgo.ComponentEmoji{
													Name: "🌐",
												},
											},
											{
												Label:       "Desktop",
												Description: "Desktop applications",
												Value:       "desktop",
												Emoji: &discordgo.ComponentEmoji{
													Name: "💻",
												},
											},
										},
									},
								},
							},
						},
					},
				}
			case "auto-populated":
				response = &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "The tastiest things are left for the end. Meet auto populated select menus.\n" +
							"By setting `MenuType` on the select menu you can tell Discord to automatically populate the menu with entities of your choice: roles, members, channels. Try one below.",
						Flags: discordgo.MessageFlagsEphemeral,
						Components: []discordgo.MessageComponent{
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.SelectMenu{
										MenuType:     discordgo.ChannelSelectMenu,
										CustomID:     "channel_select",
										Placeholder:  "Pick your favorite channel!",
										ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
									},
								},
							},
						},
					},
				}
			}
			_, err := s.InteractionRespond(i.Interaction, response)
			if err != nil {
				panic(err)
			}
		},
	}
)

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	// Components are part of interactions, so we register InteractionCreate handler
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:

			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})
	_, err := s.ApplicationCommandCreate(*AppID, *GuildID, &discordgo.ApplicationCommand{
		Name:        "buttons",
		Description: "Test the buttons if you got courage",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}
	_, err = s.ApplicationCommandCreate(*AppID, *GuildID, &discordgo.ApplicationCommand{
		Name: "selects",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "multi",
				Description: "Multi-item select menu",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "single",
				Description: "Single-item select menu",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "auto-populated",
				Description: "Automatically populated select menu, which lets you pick a member, channel or role",
			},
		},
		Description: "Lo and behold: dropdowns are coming",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}
