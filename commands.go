package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const commandPrefix = "catan!"
const helpMessage = "The available commands are: adduser, addwin, leaderboard"

type CatanBot struct {
	session        *discordgo.Session
	discordMessage *discordgo.MessageCreate
	messageParts   []string
	messageSender  MessageSender
	messageParser  MessageParser
	db             DataLayer
}

func (bot *CatanBot) handleCommand() {
	if bot.messageParser.IsCommand() {
		bot.messageParts = strings.Split(bot.discordMessage.Content, " ")
		if bot.messageParser.MessageLength() == 1 {
			bot.messageSender.sendMessage(helpMessage)
		} else if bot.messageParts[1] == "adduser" {
			bot.addUserCommand()
		} else if bot.messageParts[1] == "addwin" {
			bot.addWinCommand()
		} else if bot.messageParts[1] == "leaderboard" {
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
		row := dbConn.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE username = ($1) AND guild_id = ($2);", bot.messageParts[2], bot.discordMessage.GuildID)
		var recordExists int
		err := row.Scan(&recordExists)
		if err != nil {
			bot.session.ChannelMessageSend(bot.discordMessage.ChannelID, "An error has occurred")
			fmt.Println("Error: ", err)
			return
		}

		var response string
		if recordExists == 0 {
			response = fmt.Sprintf("User %s does not exist", bot.messageParts[2])
			bot.session.ChannelMessageSend(bot.discordMessage.ChannelID, response)
		} else {
			_, err = dbConn.Exec(context.Background(), "UPDATE users SET games_won = games_won + 1 WHERE username = ($1) AND guild_id = ($2)", bot.messageParts[2], bot.discordMessage.GuildID)
			if err != nil {
				bot.session.ChannelMessageSend(bot.discordMessage.ChannelID, "An error has occurred")
				fmt.Println("Error: ", err)
				return
			}

			messageEmbed := bot.createLeaderboardResponse()
			messageEmbed.Title = fmt.Sprintf("Congrats %s on the win!", bot.messageParts[2])

			bot.session.ChannelMessageSendEmbed(bot.discordMessage.ChannelID, messageEmbed)
		}
	} else {
		bot.session.ChannelMessageSend(bot.discordMessage.ChannelID, "Command format: addwin [username]")
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
