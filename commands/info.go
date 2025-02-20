package commands

import (
	"runtime"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command func(session *discordgo.Session, msg *discordgo.Message, args []string)

func GetCommands() map[string]Command {
	return map[string]Command{
		"info": InfoCommand,
	}
}

func InfoCommand(session *discordgo.Session, msg *discordgo.Message, args []string) {
	dependencies := []string{
		"github.com/bwmarrin/discordgo",
		"github.com/sandertv/gophertunnel/query",
	}

	embed := discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "Bot Bilgisi",
		Description: "Bu botun teknik bilgileri aşağıda verilmiştir.",
		Color:       0x3498db,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Go Sürümü",
				Value:  runtime.Version(),
				Inline: true,
			},
			{
				Name:   "Kütüphaneler",
				Value:  "```" + strings.Join(dependencies, "\n") + "```",
				Inline: false,
			},
			{
				Name:   "Author",
				Value:  "[ayd1ndemirci](https://github.com/ayd1ndemirci)",
				Inline: false,
			},
		},
	}

	_, _ = session.ChannelMessageSendEmbed(msg.ChannelID, &embed)
}
