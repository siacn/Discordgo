package discordgo_test

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func ExampleApplication() {

	// Authentication Token pulled from environment variable DGU_TOKEN
	Token := os.Getenv("DGU_TOKEN")
	if Token == "" {
		return
	}

	// Create a new Discordgo session
	dg, err := discordgo.New(Token)
	if err != nil {
		log.Println(err)
		return
	}

	// Create an new Application
	ap := &discordgo.Application{}
	ap.Name = "TestApp"
	ap.Description = "TestDesc"
	ap, err = dg.ApplicationCreate(ap)
	log.Printf("ApplicationCreate: err: %+v, app: %+v\n", err, ap)

	// Get a specific Application by it's ID
	ap, err = dg.Application(ap.ID)
	log.Printf("Application: err: %+v, app: %+v\n", err, ap)

	// Update an existing Application with new values
	ap.Description = "Whooooa"
	ap, err = dg.ApplicationUpdate(ap.ID, ap)
	log.Printf("ApplicationUpdate: err: %+v, app: %+v\n", err, ap)

	// create a new bot account for this application
	bot, err := dg.ApplicationBotCreate(ap.ID)
	log.Printf("BotCreate: err: %+v, bot: %+v\n", err, bot)

	// Get a list of all applications for the authenticated user
	apps, err := dg.Applications()
	log.Printf("Applications: err: %+v, apps : %+v\n", err, apps)
	for k, v := range apps {
		log.Printf("Applications: %d : %+v\n", k, v)
	}

	// Delete the application we created.
	err = dg.ApplicationDelete(ap.ID)
	log.Printf("Delete: err: %+v\n", err)
}
