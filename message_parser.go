package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type MessageParser interface {
	IsCommand() bool
	MessageLength() int
	GetCommand() string
}

type DiscordMessageParser struct {
	discordMessage *discordgo.MessageCreate
	parsedMessage  []string
}

func (parser *DiscordMessageParser) IsCommand() bool {
	return strings.HasPrefix(parser.discordMessage.Content, commandPrefix)
}

func (parser *DiscordMessageParser) MessageLength() int {
	if parser.parsedMessage == nil {
		parser.parsedMessage = strings.Split(parser.discordMessage.Content, " ")
	}

	return len(parser.parsedMessage)
}

func (parser *DiscordMessageParser) GetCommand() string {
	if parser.MessageLength() == 1 {
		return "" // make this return an error in the future
	} else {
		return parser.parsedMessage[1]
	}
}
