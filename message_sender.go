package main

import "github.com/bwmarrin/discordgo"

type MessageSender interface {
	sendMessage(message string)
	sendEmbedMessage(embed *discordgo.MessageEmbed)
}

type DiscordMessageSender struct {
	session        *discordgo.Session
	discordMessage *discordgo.MessageCreate
}

func (sender *DiscordMessageSender) sendMessage(message string) {
	sender.session.ChannelMessageSend(sender.discordMessage.ChannelID, message)
}

func (sender *DiscordMessageSender) sendEmbedMessage(embed *discordgo.MessageEmbed) {
	sender.session.ChannelMessageSendEmbed(sender.discordMessage.ChannelID, embed)
}
