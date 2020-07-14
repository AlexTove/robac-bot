package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"io/ioutil"
	"encoding/json"
)


func main() {

	// Read the configuration from files.
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Could not open config file,", err)
		return
	}

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Could not read config file,", err)
		return
	}

	// Parse JSON and get the bot token.
	var configs map[string]interface{}
	json.Unmarshal([]byte(bytes), &configs)

	bot_token := configs["bot_token"].(string)

	jsonFile.Close()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + bot_token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	_ = dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore all empty messages
	if len(m.Content) == 0 {
		return
	}

	if m.Content[0] == '$' {
		var (
			cmd string
		)

		r := strings.NewReader(m.Content)
		_, _ = fmt.Fscanf(r, "$%s", &cmd)

		processCommand(s, m, r, &cmd)
	}
}

func processCommand(s *discordgo.Session, m *discordgo.MessageCreate, r *strings.Reader, cmd *string) {
	switch *cmd {
	case "Baciu":
		_, _ = s.ChannelMessageSend(m.ChannelID, "E cam gay")
	case "kick":
		userPermission, _ := s.UserChannelPermissions(m.Author.ID, m.ChannelID)

		if userPermission & discordgo.PermissionBanMembers == 0 {
			return
		}

		var userid string

		_, _ = fmt.Fscanf(r, " <@!%s>", &userid)

		if len(userid) == 0 {
			return
		}

		userid = string([]rune(userid)[:len(userid) - 1])

		fmt.Println(userid)

		privateChannel, _ := s.UserChannelCreate(userid)
		_, _ = s.ChannelMessageSend(privateChannel.ID, "Sugi pula, Baciule")
		_ = s.GuildMemberDelete(m.GuildID, userid)
		s.ChannelMessageSend(m.ChannelID, "(" + m.Author.Username + ")" + " Good job! You kicked a member!")
	}
}
