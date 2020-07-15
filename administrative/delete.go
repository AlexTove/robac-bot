package administrative

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func BulkDelete(s *discordgo.Session, m *discordgo.MessageCreate, r *strings.Reader) {

	var msgCount int

	fmt.Fscanf(r, " %d", &msgCount)

	msgCount++

	messages, err := s.ChannelMessages(m.ChannelID, msgCount, "", "", "")

	if err != nil {
		fmt.Println("CH MSG")
		log.Println(err)

	}
	messagesID := make([]string, msgCount)

	if len(messages) < msgCount {
		msgCount = len(messages)
	}

	for i := 0; i < msgCount; i++ {
		messagesID[i] = messages[i].ID
		fmt.Println(messages[i].ID)
	}

	err = s.ChannelMessagesBulkDelete(m.ChannelID, messagesID)

	if err != nil {
		fmt.Println("BLK DEL MSG")
		log.Println(err)

	}
}
