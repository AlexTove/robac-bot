package administrative

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"log"
	"io"
	"strings"
)

func KickCommand(s *discordgo.Session, m *discordgo.MessageCreate, r *strings.Reader) {
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
		// TODO: add DB info here
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
