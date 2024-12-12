package events

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	commands "github.com/Basemint-Community/Confession/Commands"
	"github.com/Basemint-Community/Confession/Slash"
	"github.com/bwmarrin/discordgo"
)

func ConfessionBot() {
	token := os.Getenv("TOKEN")
	applicationID := os.Getenv("ApplicationID")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(commands.MessageCreate)
	dg.AddHandler(commands.InteractionCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		return
	}

	err = Slash.RegisterSlashCommand(dg, applicationID)
	if err != nil {
		fmt.Println("Error registering slash command: ", err)
		dg.Close()
		return
	}

	fmt.Println("Bot is now running. Press Ctrl+C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}
