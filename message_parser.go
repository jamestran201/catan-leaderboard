package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const commandPrefix = "catan!"
const minCommandLength = 2

type messageParser interface {
	isCommand() bool
	isCommandAction() bool
	numArgumentsAtLeast(int) bool
	messageLength() int
	getCommand() string
	getCommandArgument(int) string
	getGuildID() string
}

type discordMessageParser struct {
	discordMessage *discordgo.MessageCreate
	parsedMessage  []string
}

func (parser *discordMessageParser) isCommand() bool {
	return strings.HasPrefix(parser.discordMessage.Content, commandPrefix)
}

func (parser *discordMessageParser) isCommandAction() bool {
	return parser.messageLength() > 1
}

func (parser *discordMessageParser) numArgumentsAtLeast(n int) bool {
	return parser.messageLength() >= (minCommandLength + n)
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

func (parser *discordMessageParser) getCommandArgument(position int) string {
	i := minCommandLength + position - 1
	if (parser.messageLength()-1) < i || i < 1 {
		return "" // make this return an error in the future
	}

	return parser.parsedMessage[i]
}

func (parser *discordMessageParser) getGuildID() string {
	return parser.discordMessage.GuildID
}
