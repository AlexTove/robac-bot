package administrative

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func BulkDelete(s *discordgo.Session, m *discordgo.MessageCreate, r *strings.Reader) {

	var msgCount int

	_, _ = fmt.Fscanf(r, " %d", &msgCount)

	messages, _ := s.ChannelMessages(m.ChannelID, msgCount, "", "", "")

	messagesID := make([]string, msgCount)

	if len(messages) < msgCount {
		msgCount = len(messages)
	}

	for i := 0; i < msgCount; i++ {
		messagesID[i] = messages[i].ID
	}

	_ = s.ChannelMessagesBulkDelete(m.ChannelID, messagesID)

}
