package discord

import (
	"dst-admin-go/service"
	"dst-admin-go/utils/dstConfigUtils"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const Token = "MTE0MDgyMjMxNzU0NzY1NTE4OA.Gnk2Nc.X7prN532gJ1MpMtzIHnulPewe9YoNVP-Omk0m8"

var gameArchive = service.GameArchive{}
var playerService = service.PlayerService{}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "ping",
		},
		{
			Name:        "search_home",
			Description: "search home",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "name",
					Description: "home name",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		{
			Name:        "players",
			Description: "Show a list of players in the room",
		},
		{
			Name:        "baseinfo",
			Description: "Show basic information about the room",
		},
		{
			Name:        "homeinfo",
			Description: "Show basic information about the room",
		},
		{
			Name:        "respawn",
			Description: "Respawn from ghost",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "kuid",
					Description: "player kuid",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		{
			Name:        "death",
			Description: "Kill a player",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "kuid",
					Description: "player kuid",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		{
			Name:        "despawn",
			Description: "Let the player re-select the character",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "kuid",
					Description: "player kuid",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		{
			Name:        "kick",
			Description: "Kick a player",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "kuid",
					Description: "player kuid",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		{
			Name:        "dropeverything",
			Description: "Dropping items from the player's inventory",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "kuid",
					Description: "player kuid",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		{
			Name:        "rollback",
			Description: "Rollback the server a certain number of times",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "day",
					Description: "day",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		{
			Name:        "regenerateworld",
			Description: "Regenerate world",
		},
		{
			Name:        "save",
			Description: "c_save()",
		},
		{
			Name:        "shutdown",
			Description: "shutdown game",
		},
		{
			Name:        "announce",
			Description: "Send a server announcement",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "announcement",
					Description: "announcement",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
		},
		"search_home": searchHomeHandle,
		"player":      playerHandle,
		"baseinfo":    homeInfoHandle,
		"homeinfo":    homeInfoHandle,
		"respawn":     respawnHandle,
	}
)

type discordClient struct {
	token string
}

var DiscordClient *discordClient

func NewDiscordClient(token string) *discordClient {
	return &discordClient{
		token: token,
	}
}

func (d *discordClient) Start() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + d.token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	registeredCommand(dg)

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "/ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "/pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func searchHomeHandle(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	if option, ok := optionMap["name"]; ok {
		name := option.StringValue()
		log.Println(name)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Search Home " + "\"" + name + "\"",
						Color:       0x5c59bd,
						Type:        discordgo.EmbedTypeRich,
						Description: "find 8 homes，page 1 total page 1",
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "猜猜我是谁的世界1(1/8)",
								Value: "\u200B",
							},
							{
								Name:  "猜猜我是谁的世界2(2/10)",
								Value: "\u200B",
							},
						},
					},
				},
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Search error, Please input home name",
		},
	})
}

func playerHandle(s *discordgo.Session, i *discordgo.InteractionCreate) {

	config := dstConfigUtils.GetDstConfig()
	clusterName := config.Cluster
	playerlist := playerService.GetPlayerList(clusterName, "Master")

	var nameset string
	var roleset string
	var dayset string
	var kuidset string
	for i := range playerlist {
		nameset = nameset + playerlist[i].Name + "(" + playerlist[i].KuId + ")" + "\n"
		roleset = roleset + playerlist[i].Role + "\n"
		dayset = dayset + playerlist[i].Day + " day" + "\n"
		kuidset = kuidset + playerlist[i].KuId + " " + "\n"

	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title: "Player List",
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Player List",
					Color: 0x376dd9,
					Type:  discordgo.EmbedTypeRich,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Name",
							Value:  nameset,
							Inline: true,
						},
						{
							Name:   "Role",
							Value:  roleset,
							Inline: true,
						},
						{
							Name:   "Day",
							Value:  dayset,
							Inline: true,
						},
						{
							Name:   "KuId",
							Value:  kuidset,
							Inline: true,
						},
					},
				},
			},
		},
	})

	if err != nil {
		return
	}
}

func homeInfoHandle(s *discordgo.Session, i *discordgo.InteractionCreate) {

	config := dstConfigUtils.GetDstConfig()
	clusterName := config.Cluster
	archive := gameArchive.GetGameArchive(clusterName)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title: "Home archive info",
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       archive.ClusterName,
					Description: archive.ClusterDescription,
					Color:       0x21dd7f,
					Type:        discordgo.EmbedTypeRich,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Players",
							Value: strconv.Itoa(len(archive.Players)) + "/" + strconv.Itoa(archive.MaxPlayers),
						},
						{
							Name:  "GameMod",
							Value: archive.GameMod,
						},
						{
							Name:  "Mods",
							Value: strconv.Itoa(archive.Mods) + "nums",
						},
						{
							Name:  "Day",
							Value: strconv.Itoa(archive.Meta.Clock.Cycles) + "days",
						},
						{
							Name:  "Season",
							Value: archive.Meta.Seasons.Season + "(" + strconv.Itoa(archive.Meta.Seasons.ElapsedDaysInSeason) + "/" + strconv.Itoa(archive.Meta.Seasons.ElapsedDaysInSeason+archive.Meta.Seasons.RemainingDaysInSeason) + ")",
						},
						{
							Name:  "Ip Connect",
							Value: archive.IpConnect,
						},
						{
							Name:  "Version",
							Value: strconv.FormatInt(archive.Version, 10) + "/" + strconv.FormatInt(archive.LastVersion, 10),
						},
					},
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "Update Game",
							Style: discordgo.PrimaryButton,
							URL:   "https://discord.com/developers/docs/interactions/message-components#buttons1",
						},
						discordgo.Button{
							Label: "Backup Game",
							Style: discordgo.PrimaryButton,
							URL:   "https://discord.com/developers/docs/interactions/message-components#buttons2",
						},
					},
				},
			},
		},
	})

	if err != nil {
		return
	}
}

func respawnHandle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// Get the value from the option map.
	// When the option exists, ok = true
	if option, ok := optionMap["kuid"]; ok {
		kuid := option.StringValue()
		log.Println(kuid)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, they will be discussed in "responses"
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "复活" + kuid + "玩家成功 ",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "please input ku_id, ku_id must not be null",
		},
	})
}

func registeredCommand(s *discordgo.Session) {

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}
