package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const commandPrefix = "catan!"
const helpMessage = "The available commands are: adduser, addwin, leaderboard"

type catanBot struct {
	messageSender messageSender
	messageParser messageParser
	db            dataLayer
}

func (bot *catanBot) handleCommand() {
	if !bot.messageParser.isCommand() {
		return
	}

	if !bot.messageParser.isCommandAction() {
		bot.messageSender.sendMessage(helpMessage)
		return
	}

	switch bot.messageParser.getCommand() {
	case "adduser":
		bot.addUser()
	case "addwin":
		bot.addWin()
	case "leaderboard":
		bot.showLeaderboard()
	default:
		bot.messageSender.sendMessage(helpMessage)
	}
}

func (bot *catanBot) addUser() {
	if !bot.messageParser.isCommandWithArgument() {
		bot.messageSender.sendMessage("Command format: adduser [username]")
		return
	}

	err := bot.db.addUser(bot.messageParser.getCommandArgument(), bot.messageParser.getGuildID())
	if err != nil {
		bot.messageSender.sendMessage("An error has occurred")
		fmt.Println("Error: ", err)
		return
	}

	response := fmt.Sprintf("Successfully added user: %s", bot.messageParser.getCommandArgument())
	bot.messageSender.sendMessage(response)
}

func (bot *catanBot) addWin() {
	if !bot.messageParser.isCommandWithArgument() {
		bot.messageSender.sendMessage("Command format: addwin [username]")
		return
	}

	username := bot.messageParser.getCommandArgument()
	guildID := bot.messageParser.getGuildID()

	recordExists, err := bot.db.checkUserExists(username, guildID)

	if err != nil {
		bot.messageSender.sendMessage("An error has occurred")
		fmt.Println("Error: ", err)
		return
	}

	if recordExists == 0 {
		response := fmt.Sprintf("User %s does not exist", username)
		bot.messageSender.sendMessage(response)
	} else {
		err = bot.db.addWin(username, guildID)

		if err != nil {
			bot.messageSender.sendMessage("An error has occurred")
			fmt.Println("Error: ", err)
			return
		}

		messageEmbed := bot.createLeaderboardResponse()
		messageEmbed.Title = fmt.Sprintf("Congrats %s on the win!", username)

		bot.messageSender.sendEmbedMessage(messageEmbed)
	}
}

func (bot *catanBot) showLeaderboard() {
	bot.messageSender.sendEmbedMessage(bot.createLeaderboardResponse())
}

func (bot *catanBot) createLeaderboardResponse() *discordgo.MessageEmbed {
	users, err := bot.db.getTopFiveUsers(bot.messageParser.getGuildID())

	if err != nil {
		bot.messageSender.sendMessage("An error has occurred")
		fmt.Println("Error: ", err)
		return nil
	}

	var (
		rankField      = discordgo.MessageEmbedField{"Rank", "", true}
		usernameField  = discordgo.MessageEmbedField{"Username", "", true}
		victoriesField = discordgo.MessageEmbedField{"Victories", "", true}
		ranks          = make([]string, 0, 5)
		usernames      = make([]string, 0, 5)
		victories      = make([]string, 0, 5)
	)

	for _, user := range users {
		ranks = append(ranks, user.rank)
		usernames = append(usernames, user.username)
		victories = append(victories, user.victories)
	}

	rankField.Value = strings.Join(ranks, "\n")
	usernameField.Value = strings.Join(usernames, "\n")
	victoriesField.Value = strings.Join(victories, "\n")

	messageEmbed := discordgo.MessageEmbed{}
	messageEmbed.Fields = []*discordgo.MessageEmbedField{&rankField, &usernameField, &victoriesField}

	return &messageEmbed
}
