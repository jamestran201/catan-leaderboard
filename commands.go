package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const commandPrefix = "catan!"
const helpMessage = "The available commands are: adduser, addwin, leaderboard"

type CatanBot struct {
	messageSender MessageSender
	messageParser MessageParser
	db            DataLayer
}

func (bot *CatanBot) handleCommand() {
	if bot.messageParser.IsCommand() {
		if bot.messageParser.MessageLength() == 1 {
			bot.messageSender.sendMessage(helpMessage)
		} else if bot.messageParser.GetCommand() == "adduser" {
			bot.addUserCommand()
		} else if bot.messageParser.GetCommand() == "addwin" {
			bot.addWinCommand()
		} else if bot.messageParser.GetCommand() == "leaderboard" {
			bot.showLeaderboardCommand()
		} else {
			bot.messageSender.sendMessage(helpMessage)
		}
	}
}

func (bot *CatanBot) addUserCommand() {
	if bot.messageParser.MessageLength() == 3 {
		err := bot.db.AddUser(bot.messageParser.GetCommandArgument(), bot.messageParser.GetGuildID())
		if err != nil {
			bot.messageSender.sendMessage("An error has occurred")
			fmt.Println("Error: ", err)
			return
		}

		response := fmt.Sprintf("Successfully added user: %s", bot.messageParser.GetCommandArgument())
		bot.messageSender.sendMessage(response)
	} else {
		bot.messageSender.sendMessage("Command format: adduser [username]")
	}
}

func (bot *CatanBot) addWinCommand() {
	if bot.messageParser.MessageLength() == 3 {
		username := bot.messageParser.GetCommandArgument()
		guildID := bot.messageParser.GetGuildID()

		recordExists, err := bot.db.CheckUserExists(username, guildID)

		if err != nil {
			bot.messageSender.sendMessage("An error has occurred")
			fmt.Println("Error: ", err)
			return
		}

		if recordExists == 0 {
			response := fmt.Sprintf("User %s does not exist", username)
			bot.messageSender.sendMessage(response)
		} else {
			err = bot.db.AddWin(username, guildID)

			if err != nil {
				bot.messageSender.sendMessage("An error has occurred")
				fmt.Println("Error: ", err)
				return
			}

			messageEmbed := bot.createLeaderboardResponse()
			messageEmbed.Title = fmt.Sprintf("Congrats %s on the win!", username)

			bot.messageSender.sendEmbedMessage(messageEmbed)
		}
	} else {
		bot.messageSender.sendMessage("Command format: addwin [username]")
	}
}

func (bot *CatanBot) showLeaderboardCommand() {
	bot.messageSender.sendEmbedMessage(bot.createLeaderboardResponse())
}

func (bot *CatanBot) createLeaderboardResponse() *discordgo.MessageEmbed {
	users, err := bot.db.GetTopFiveUsers(bot.messageParser.GetGuildID())

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
		ranks = append(ranks, user.Rank)
		usernames = append(usernames, user.Username)
		victories = append(victories, user.Victories)
	}

	rankField.Value = strings.Join(ranks, "\n")
	usernameField.Value = strings.Join(usernames, "\n")
	victoriesField.Value = strings.Join(victories, "\n")

	messageEmbed := discordgo.MessageEmbed{}
	messageEmbed.Fields = []*discordgo.MessageEmbedField{&rankField, &usernameField, &victoriesField}

	return &messageEmbed
}
