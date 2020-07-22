package main

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestSendHelpMessageForUnknownCommand(t *testing.T) {
	message := &discordgo.Message{Content: "catan! random"}
	messageCreate := &discordgo.MessageCreate{message}
	sender := MessageSenderMock{}
	bot := &CatanBot{nil, messageCreate, nil, &sender}

	bot.handleCommand()

	if sender.messageSent != helpMessage {
		t.Errorf("Got %s\nWant %s", sender.messageSent, helpMessage)
	}
}

type MessageSenderMock struct {
	messageSent string
}

func (sender *MessageSenderMock) sendMessage(message string) {
	sender.messageSent = message
}