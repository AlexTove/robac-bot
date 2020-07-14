package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
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
	_ = json.Unmarshal([]byte(bytes), &configs)

	bot_token := configs["bot_token"].(string)

	_ = jsonFile.Close()

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

		if userPermission & discordgo.PermissionKickMembers == 0 ||
			userPermission & discordgo.PermissionAdministrator == 0 {
			return
		}

		var userid, reason string

		_, _ = fmt.Fscanf(r, " %s", &userid)

		reg, err := regexp.Compile("[^0-9]+")

		if err != nil {
			log.Fatal(err)
		}

		userid = reg.ReplaceAllString(userid, "")

		if len(userid) == 0 {
			return
		}

		buf := new(strings.Builder)
		_, _ = io.Copy(buf, r)

		if len(buf.String()) == 0 {
			reason = "No reason was given."
		} else {
			reason = buf.String()
		}

		err = s.GuildMemberDelete(m.GuildID, userid)

		if userid == m.Author.ID {
			_, _ = s.ChannelMessageSend(m.ChannelID, "(" + m.Author.Username + ")" +
				" You can't kick yourself!")
			return
		}

		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "(" + m.Author.Username + ")" +
				" You can't kick that user!")
			return
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, "(" + m.Author.Username + ")" +
				" User kicked successfully!")
		}

		privateChannel, err := s.UserChannelCreate(userid)

		if err != nil {
			log.Fatal(err)
		}

		_, err = s.ChannelMessageSend(privateChannel.ID, "You have been kicked from the server.\n" +
			"Reason: " + reason + "\n" +
			"Author: " + m.Author.Username)

		if err != nil {
			log.Fatal(err)
		}
	}
}
