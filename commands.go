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
}

func (bot *CatanBot) handleCommand() {
	if strings.HasPrefix(bot.discordMessage.Content, commandPrefix) {
		bot.messageParts = strings.Split(bot.discordMessage.Content, " ")
		if len(bot.messageParts) == 1 {
			bot.sendHelpMessage()
		} else if bot.messageParts[1] == "adduser" {
			bot.addUserCommand()
		} else if bot.messageParts[1] == "addwin" {
			bot.addWinCommand()
		} else if bot.messageParts[1] == "leaderboard" {
			bot.showLeaderboardCommand()
		} else {
			bot.sendHelpMessage()
		}
	}
}

func (bot *CatanBot) sendHelpMessage() {
	bot.session.ChannelMessageSend(bot.discordMessage.ChannelID, helpMessage)
}

func (bot *CatanBot) addUserCommand() {
	if len(bot.messageParts) == 3 {
		_, err := dbConn.Exec(context.Background(), "INSERT INTO users (username, guild_id) VALUES ($1, $2)", bot.messageParts[2], bot.discordMessage.GuildID)
		if err != nil {
			bot.session.ChannelMessageSend(bot.discordMessage.ChannelID, "An error has occurred")
			fmt.Println("Error: ", err)
			return
		}

		response := fmt.Sprintf("Successfully added user: %s", bot.messageParts[2])
		bot.session.ChannelMessageSend(bot.discordMessage.ChannelID, response)
	} else {
		bot.session.ChannelMessageSend(bot.discordMessage.ChannelID, "Command format: adduser [username]")
	}
}

func (bot *CatanBot) addWinCommand() {
	if len(bot.messageParts) == 3 {
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
	bot.session.ChannelMessageSendEmbed(bot.discordMessage.ChannelID, bot.createLeaderboardResponse())
}

func (bot *CatanBot) createLeaderboardResponse() *discordgo.MessageEmbed {
	rows, err := dbConn.Query(context.Background(), "SELECT CAST(RANK() OVER (ORDER BY games_won DESC) AS TEXT), username, CAST(games_won AS TEXT) FROM users WHERE guild_id = ($1) LIMIT 5", bot.discordMessage.GuildID)
	if err != nil {
		bot.session.ChannelMessageSend(bot.discordMessage.ChannelID, "An error has occurred")
		fmt.Println("Error: ", err)
		return nil
	}

	defer rows.Close()

	var (
		rankField      = discordgo.MessageEmbedField{"Rank", "", true}
		usernameField  = discordgo.MessageEmbedField{"Username", "", true}
		victoriesField = discordgo.MessageEmbedField{"Victories", "", true}

		ranks     [5]string
		usernames [5]string
		victories [5]string
	)

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&ranks[i], &usernames[i], &victories[i])
		if err != nil {
			bot.session.ChannelMessageSend(bot.discordMessage.ChannelID, "An error has occurred")
			fmt.Println("Error: ", err)
			return nil
		}
	}

	rankField.Value = strings.Join(ranks[:], "\n")
	usernameField.Value = strings.Join(usernames[:], "\n")
	victoriesField.Value = strings.Join(victories[:], "\n")

	if rows.Err() != nil {
		bot.session.ChannelMessageSend(bot.discordMessage.ChannelID, "An error has occurred")
		fmt.Println("Error: ", err)
		return nil
	}

	messageEmbed := discordgo.MessageEmbed{}
	messageEmbed.Fields = []*discordgo.MessageEmbedField{&rankField, &usernameField, &victoriesField}
	return &messageEmbed
}
