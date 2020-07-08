package main

import (
  "github.com/bwmarrin/discordgo"
  "github.com/jackc/pgx/v4"
  "context"
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "strings"
)

const COMMAND_PREFIX = "catan!"

var db_conn *pgx.Conn

func init() {
  var err error
  db_conn, err = pgx.Connect(context.Background(), os.Getenv("POSTGRESQL_URL"))
  if err != nil {
    fmt.Println("Error connecting to the database: ", err)
    os.Exit(1)
  }
}

func main() {
  defer db_conn.Close(context.Background())

  token := os.Getenv("BOT_TOKEN")
  discord, err := discordgo.New("Bot " + token)
  if err != nil {
    fmt.Println("Error creating Discord session: ", err)
    os.Exit(1)
  }

  discord.AddHandler(messageCreate)

  err = discord.Open()
  if err != nil {
    fmt.Println("Error opening connection: ", err)
    os.Exit(1)
  }

  fmt.Println("Bot is now running. Press CTRL-C to exit.")
  sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
  <-sc

  discord.Close()
}

func messageCreate(session *discordgo.Session, m *discordgo.MessageCreate) {
  if m.Author.ID == session.State.User.ID {
    return
  }

  if strings.HasPrefix(m.Content, COMMAND_PREFIX) {
    message := strings.Split(m.Content, " ")
    if message[1] == "adduser" {
      addUserCommand(session, m, message)
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
