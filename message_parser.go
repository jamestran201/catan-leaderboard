package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type messageParser interface {
	isCommand() bool
	messageLength() int
	getCommand() string
	getCommandArgument() string
	getGuildID() string
}

type discordMessageParser struct {
	discordMessage *discordgo.MessageCreate
	parsedMessage  []string
}

func (parser *discordMessageParser) isCommand() bool {
	return strings.HasPrefix(parser.discordMessage.Content, commandPrefix)
}

func (parser *discordMessageParser) messageLength() int {
	if parser.parsedMessage == nil {
		parser.parsedMessage = strings.Split(parser.discordMessage.Content, " ")
	}

	return len(parser.parsedMessage)
}

func (parser *discordMessageParser) getCommand() string {
	if parser.messageLength() == 1 {
		return "" // make this return an error in the future
	}

	return parser.parsedMessage[1]
}

func (parser *discordMessageParser) getCommandArgument() string {
	if parser.messageLength() < 3 {
		return "" // make this return an error in the future
	}

	return parser.parsedMessage[2]
}

func (parser *discordMessageParser) getGuildID() string {
	return parser.discordMessage.GuildID
}
