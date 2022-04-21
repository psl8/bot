package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
    "crypto/rand"
    "math/big"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

type BotConfig struct {
	Token string
}

func execSlashCommand(cmd *discord.CommandInteraction, s *state.State, e *gateway.InteractionCreateEvent) {
	switch cmd.Name {
	case "ping":
		data := api.InteractionResponse{
			Type: api.MessageInteractionWithSource,
			Data: &api.InteractionResponseData{
				Content: option.NewNullableString("Pong!"),
			},
		}

		if err := s.RespondInteraction(e.ID, e.Token, data); err != nil {
			log.Println("failed to send interaction callback:", err)
		}
    case "rr":
        nBig, err := rand.Int(rand.Reader, big.NewInt(1))
        if err != nil {
			log.Println("failed to get random int:", err)
            return;
        }
        response_string := ""
        if nBig.Int64() == 1 {
            response_string = "*bang*\nYou're dead :("
        } else {
            response_string = "*click*\nYou're safe \U0001f920"
        }

		data := api.InteractionResponse{
			Type: api.MessageInteractionWithSource,
			Data: &api.InteractionResponseData{
				Content: option.NewNullableString(response_string),
			},
		}

		if err := s.RespondInteraction(e.ID, e.Token, data); err != nil {
			log.Println("failed to send interaction callback:", err)
		}
	default:
		log.Printf("unknown command %#v %#v %#v", cmd, s, e)
	}
}

func main() {
	configFile, err := os.ReadFile("botconfig.json")
	if err != nil {
		log.Fatalln("failed to read configuration:", err)
	}

	var config BotConfig
	if err := json.Unmarshal(configFile, &config); err != nil {
		log.Fatalln("failed to parse configuration:", err)
	}

	s := state.New("Bot " + config.Token)

	app, err := s.CurrentApplication()
	if err != nil {
		log.Fatalln("failed to get application ID:", err)
	}

	s.AddHandler(func(e *gateway.InteractionCreateEvent) {
		switch v := e.Data.(type) {
		case *discord.CommandInteraction:
			execSlashCommand(v, s, e)
		default:
			log.Println("unknown event:", e)
		}
	})

	s.AddIntents(gateway.IntentGuilds)
	s.AddIntents(gateway.IntentGuildMessages)

	if err := s.Open(context.Background()); err != nil {
		log.Fatalln("failed to open:", err)
	}
	defer s.Close()

	log.Println("Gateway connected. Registering commands.")

	log.Println("Logging existing commands")
	commands, err := s.Commands(app.ID)
	if err != nil {
		log.Fatalln("failed to get commands:", err)
	}

	for _, command := range commands {
		log.Println("Existing command", command.Name, "found.")
	}

	newCommands := []api.CreateCommandData{
		{
			Name:        "ping",
			Description: "Basic ping command.",
		},
		{
			Name:        "rr",
            Description: "Russian roulette pew pew",
		},
	}

	if _, err := s.BulkOverwriteCommands(app.ID, newCommands); err != nil {
		log.Fatalln("failed to create command:", err)
	}

	// Block forever.
	select {}
}
