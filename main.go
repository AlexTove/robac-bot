package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

// Variables used for command line parameters
const Token string = "NzMxMzkwMjU3MjQ4ODYyMjc4.Xwoq8w.FgS6svDElhlNndt9yjQ8I0tMPW8"

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
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
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	//If the message is "Hello" reply back
	if m.Content == "Hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello, " + m.Author.Mention() + "!")
	}

	if m.Content == "Sugi pula" {
		s.ChannelMessageSend(m.ChannelID, "Te bag eu in pizda ma-tii bai " + m.Author.Mention())
	}

	if(m.Content[0] == '$') {
		strlen := len(m.Content)
		command := ""

		for i := 1; i != strlen; i++{
			if m.Content[i] != ' ' && m.Content[i] != '\n' {
				command += string(m.Content[i])
			} else {
				break
			}
		}
		s.ChannelMessageSend(m.ChannelID, command)
	}
}