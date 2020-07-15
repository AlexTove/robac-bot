package administrative

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"log"
	"io"
	"strings"
)

// Kick command
// SYNTAX:
// *kick user
//	Where '*' is the command prefix and 'user' is the mentioned user (@...)
//	Will try to kick 'user'. No reason is specified in this form, in the private message to 'user'.
// *kick user msg
//	Where 'msg' is the reason for the kick
//	Will try to kick 'user'. A given reason is put instead, in the private message to 'user'.
func KickCommand(s *discordgo.Session, m *discordgo.MessageCreate, r *strings.Reader) {
	userPermission, _ := s.UserChannelPermissions(m.Author.ID, m.ChannelID)

	// Verify the user's permissions
	if userPermission & discordgo.PermissionKickMembers == 0 ||
		userPermission & discordgo.PermissionAdministrator == 0 {
		return
	}

	// Get the userid with regex
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

	// Check for the second argument
	buf := new(strings.Builder)
	_, _ = io.Copy(buf, r)

	if len(buf.String()) == 0 {
		reason = "No reason was given."
	} else {
		reason = buf.String()
	}

	// Attempt to kick the user
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
		// TODO: add DB info here
	}

	// Send a private message to the user
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
