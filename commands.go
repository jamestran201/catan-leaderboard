package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func handleCommands(session *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Content, COMMAND_PREFIX) {
		message := strings.Split(m.Content, " ")
		help_message := "The available commands are: adduser, addwin, leaderboard"
		if len(message) == 1 {
			session.ChannelMessageSend(m.ChannelID, help_message)
		} else if message[1] == "adduser" {
			addUserCommand(session, m, message)
		} else if message[1] == "addwin" {
			addWinCommand(session, m, message)
		} else if message[1] == "leaderboard" {
			showLeaderboardCommand(session, m, message)
		} else {
			session.ChannelMessageSend(m.ChannelID, help_message)
		}
	}
}

func addUserCommand(session *discordgo.Session, m *discordgo.MessageCreate, message []string) {
	if len(message) == 3 {
		_, err := db_conn.Exec(context.Background(), "INSERT INTO users (username, guild_id) VALUES ($1, $2)", message[2], m.GuildID)
		if err != nil {
			session.ChannelMessageSend(m.ChannelID, "An error has occurred")
			fmt.Println("Error: ", err)
			return
		}

		response := fmt.Sprintf("Successfully added user: %s", message[2])
		session.ChannelMessageSend(m.ChannelID, response)
	} else {
		session.ChannelMessageSend(m.ChannelID, "Command format: adduser [username]")
	}
}

func addWinCommand(session *discordgo.Session, m *discordgo.MessageCreate, message []string) {
	if len(message) == 3 {
		row := db_conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE username = ($1) AND guild_id = ($2);", message[2], m.GuildID)
		var record_exists int
		err := row.Scan(&record_exists)
		if err != nil {
			session.ChannelMessageSend(m.ChannelID, "An error has occurred")
			fmt.Println("Error: ", err)
			return
		}

		var response string
		if record_exists == 0 {
			response = fmt.Sprintf("User %s does not exist", message[2])
		} else {
			_, err = db_conn.Exec(context.Background(), "UPDATE users SET games_won = games_won + 1 WHERE username = ($1) AND guild_id = ($2)", message[2], m.GuildID)
			if err != nil {
				session.ChannelMessageSend(m.ChannelID, "An error has occurred")
				fmt.Println("Error: ", err)
				return
			}

			response = fmt.Sprintf("Congrats %s on the win!", message[2])
		}

		session.ChannelMessageSend(m.ChannelID, response)
	} else {
		session.ChannelMessageSend(m.ChannelID, "Command format: addwin [username]")
	}
}

func showLeaderboardCommand(session *discordgo.Session, m *discordgo.MessageCreate, message []string) {
	rows, err := db_conn.Query(context.Background(), "SELECT CAST(RANK() OVER (ORDER BY games_won DESC) AS TEXT), username, CAST(games_won AS TEXT) FROM users WHERE guild_id = ($1) LIMIT 5", m.GuildID)
	if err != nil {
		session.ChannelMessageSend(m.ChannelID, "An error has occurred")
		fmt.Println("Error: ", err)
		return
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
			session.ChannelMessageSend(m.ChannelID, "An error has occurred")
			fmt.Println("Error: ", err)
			return
		}
	}

	rankField.Value = strings.Join(ranks[:], "\n")
	usernameField.Value = strings.Join(usernames[:], "\n")
	victoriesField.Value = strings.Join(victories[:], "\n")

	if rows.Err() != nil {
		session.ChannelMessageSend(m.ChannelID, "An error has occurred")
		fmt.Println("Error: ", err)
		return
	}

	messageEmbed := discordgo.MessageEmbed{}
	messageEmbed.Fields = []*discordgo.MessageEmbedField{&rankField, &usernameField, &victoriesField}
	session.ChannelMessageSendEmbed(m.ChannelID, &messageEmbed)
}
