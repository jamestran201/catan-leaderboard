package main

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestSendHelpMessageForUnknownCommand(t *testing.T) {
	message := &discordgo.Message{Content: "catan! random"}
	messageCreate := &discordgo.MessageCreate{message}
	sender := MessageSenderMock{}
	messageParser := DiscordMessageParser{discordMessage: messageCreate}
	bot := &CatanBot{nil, messageCreate, nil, &sender, &messageParser, nil}

	bot.handleCommand()

	if sender.messageSent != helpMessage {
		t.Errorf("Got %s\nWant %s", sender.messageSent, helpMessage)
	}
}

func TestSendHelpMessage(t *testing.T) {
	message := &discordgo.Message{Content: "catan!"}
	messageCreate := &discordgo.MessageCreate{message}
	sender := MessageSenderMock{}
	messageParser := DiscordMessageParser{discordMessage: messageCreate}
	bot := &CatanBot{nil, messageCreate, nil, &sender, &messageParser, nil}

	bot.handleCommand()

	if sender.messageSent != helpMessage {
		t.Errorf("Got %s\nWant %s", sender.messageSent, helpMessage)
	}
}

func TestAddUserSuccess(t *testing.T) {
	message := &discordgo.Message{Content: "catan! adduser test_user"}
	messageCreate := &discordgo.MessageCreate{message}
	sender := MessageSenderMock{}
	messageParser := DiscordMessageParser{discordMessage: messageCreate}
	db := &MockDataLayer{}
	bot := &CatanBot{nil, messageCreate, nil, &sender, &messageParser, db}

	bot.handleCommand()

	expected := "Successfully added user: test_user"
	if sender.messageSent != expected {
		t.Errorf("Got %s\nWant %s", sender.messageSent, expected)
	}
}

func TestAddUserWrongFormat(t *testing.T) {
	message := &discordgo.Message{Content: "catan! adduser"}
	messageCreate := &discordgo.MessageCreate{message}
	sender := MessageSenderMock{}
	messageParser := DiscordMessageParser{discordMessage: messageCreate}
	db := &MockDataLayer{}
	bot := &CatanBot{nil, messageCreate, nil, &sender, &messageParser, db}

	bot.handleCommand()

	expected := "Command format: adduser [username]"
	if sender.messageSent != expected {
		t.Errorf("Got %s\nWant %s", sender.messageSent, expected)
	}
}

type MessageSenderMock struct {
	messageSent string
}

func (sender *MessageSenderMock) sendMessage(message string) {
	sender.messageSent = message
}

func (sender *MessageSenderMock) sendEmbedMessage(embed *discordgo.MessageEmbed) {}

type MockDataLayer struct{}

func (db *MockDataLayer) AddUser(username string, guild_id string) error {
	return nil
}

func (db *MockDataLayer) GetTopFiveUsers(guildID string) ([]User, error) {
	return nil, nil
}
