package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"os/signal"
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
	_ = json.Unmarshal(bytes, &configs)

	botToken := configs["botToken"].(string)

	_ = jsonFile.Close()

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

	botToken := configs["botToken"].(string)

	_ = jsonFile.Close()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + botToken)
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

	// Check for the bot-specific command prefix
	if m.Content[0] == '$' {
		var (
			cmd string
		)

		r := strings.NewReader(m.Content)
		_, _ = fmt.Fscanf(r, "$%s", &cmd)

		processCommand(s, m, r, &cmd)
	}
}
