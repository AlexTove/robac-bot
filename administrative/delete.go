package administrative

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func bulkDelete(s *discordgo.Session, m *discordgo.MessageCreate, r *strings.Reader) {

	var msgCount int

	fmt.Fscanf(r, " %d", &msgCount)

	messages, _ := s.ChannelMessages(m.ChannelID, msgCount, "", "", "")

	messagesID := make([]string, msgCount)

	if len(messages) < msgCount {
		msgCount = len(messages)
	}

	for i := 0; i < msgCount; i++ {
		messagesID[i] = messages[i].ID
	}

	s.ChannelMessagesBulkDelete(m.ChannelID, messagesID)

}
