package main

import (
	"strings"

	"github.com/adriangeorge/robac-bot/administrative"
	"github.com/bwmarrin/discordgo"
)

// The route function for commands, from where the command functions are executed.
func processCommand(s *discordgo.Session, m *discordgo.MessageCreate, r *strings.Reader, cmd *string) {
	switch *cmd {
	case "Baciu":
		_, _ = s.ChannelMessageSend(m.ChannelID, "E cam gay")

	case "kick":
		administrative.KickCommand(s, m, r)

	case "delete":
		administrative.BulkDelete(s, m, r)
	}

}
