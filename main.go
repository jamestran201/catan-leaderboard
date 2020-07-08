package main

import (
  "github.com/bwmarrin/discordgo"
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "strings"
)

const COMMAND_PREFIX = "catanleaderboard!"

func main() {
  fmt.Println("Testing Discord bot")

  token := os.Getenv("BOT_TOKEN")
  discord, err := discordgo.New("Bot " + token)
  if err != nil {
    fmt.Println("Error creating Discord session: ", err)
    return
  }

  discord.AddHandler(messageCreate)

  err = discord.Open()
  if err != nil {
    fmt.Println("Error opening connection: ", err)
    return
  }

  fmt.Println("Bot is now running. Press CTRL-C to exit.")
  sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
  <-sc

  discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
  if m.Author.ID == s.State.User.ID {
    return
  }

  if strings.HasPrefix(m.Content, COMMAND_PREFIX) {
    s.ChannelMessageSend(m.ChannelID, "Command received!")
  }
}
