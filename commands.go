package main

import (
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
)

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
	case "record":
		bot.recordGame()
	default:
		bot.messageSender.sendMessage(helpMessage)
	}
}

func (bot *catanBot) addUser() {
	if !bot.messageParser.numArgumentsAtLeast(1) {
		bot.messageSender.sendMessage("Command format: adduser [username]")
		return
	}

	username := bot.messageParser.getCommandArgument(1)
	err := bot.db.addUser(username, bot.messageParser.getGuildID())
	if err != nil {
		bot.messageSender.sendMessage("An error has occurred")
		fmt.Println("Error: ", err)
		return
	}

	response := fmt.Sprintf("Successfully added user: %s", username)
	bot.messageSender.sendMessage(response)
}

func (bot *catanBot) addWin() {
	if !bot.messageParser.numArgumentsAtLeast(1) {
		bot.messageSender.sendMessage("Command format: addwin [username]")
		return
	}

	username := bot.messageParser.getCommandArgument(1)
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

		congratsMessage := fmt.Sprintf("Congrats %s on the win! :tada:\n", username)

		leaderboardMessage, err := bot.createLeaderboardResponse()

		if err != nil {
			bot.messageSender.sendMessage("An error has occurred")
			fmt.Println("Error: ", err)
			return
		}

		message := fmt.Sprintf("%s%s", congratsMessage, leaderboardMessage)
		bot.messageSender.sendMessage(message)
	}
}

func (bot *catanBot) showLeaderboard() {
	message, err := bot.createLeaderboardResponse()

	if err != nil {
		bot.messageSender.sendMessage("An error has occurred")
		fmt.Println("Error: ", err)
		return
	}

	bot.messageSender.sendMessage(message)
}

func (bot *catanBot) recordGame() {

}

func (bot *catanBot) createLeaderboardResponse() (string, error) {
	users, err := bot.db.getTopTwentyUsers(bot.messageParser.getGuildID())

	if err != nil {
		return "", err
	}

	var stringBuilder strings.Builder
	table := tablewriter.NewWriter(&stringBuilder)
	table.SetHeader([]string{"Rank", "Username", "Victories", "Points", "Games"})

	for _, user := range users {
		data := []string{user.rank, user.username, user.victories, user.points, user.games}
		table.Append(data)
	}

	table.Render()

	return fmt.Sprintf("```%s```", stringBuilder.String()), nil
}
