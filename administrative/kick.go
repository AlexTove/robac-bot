package administrative

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"regexp"
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
	if userPermission&discordgo.PermissionKickMembers == 0 &&
		userPermission&discordgo.PermissionAdministrator == 0 {
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

	if userid == m.Author.ID {
		_, _ = s.ChannelMessageSend(m.ChannelID, "("+m.Author.Username+")"+
			" You can't kick yourself!")
		return
	}

	_, err = s.GuildMember(m.GuildID, userid)

	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "("+m.Author.Username+")"+
			" You must provide a valid userID!")
		return
	}

	kickedUserPermission, _ := s.UserChannelPermissions(userid, m.ChannelID)

	if kickedUserPermission&discordgo.PermissionKickMembers != 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "("+m.Author.Username+")"+
			" You can't kick that user!")
		return
	} else {
		// Send a private message to the user
		privateChannel, _ := s.UserChannelCreate(userid)

		_, err = s.ChannelMessageSend(privateChannel.ID, "You have been kicked from the server.\n"+
			"Reason: "+reason+"\n"+
			"Author: "+m.Author.Username)

		_, _ = s.ChannelMessageSend(m.ChannelID, "("+m.Author.Username+")"+
			" User kicked successfully!")

		_ = s.GuildMemberDelete(m.GuildID, userid)

		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/robacbot")

		if err != nil {
			panic(err)
		}

		defer db.Close()

		query := "INSERT INTO kicklog (kickedUserID, kickedByUserID, reason) VALUES (?, ?, ?)"
		insert, err := db.Query(query, userid, m.Author.ID, reason)

		if err != nil {
			panic(err)
		}

		insert.Close()
	}

	if err != nil {
		log.Fatal(err)
	}

}
