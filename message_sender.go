package main

import "github.com/bwmarrin/discordgo"

type messageSender interface {
	sendMessage(message string)
	sendEmbedMessage(embed *discordgo.MessageEmbed)
}

type discordMessageSender struct {
	session        *discordgo.Session
	discordMessage *discordgo.MessageCreate
}

func (sender *discordMessageSender) sendMessage(message string) {
	sender.session.ChannelMessageSend(sender.discordMessage.ChannelID, message)
}

func (sender *discordMessageSender) sendEmbedMessage(embed *discordgo.MessageEmbed) {
	sender.session.ChannelMessageSendEmbed(sender.discordMessage.ChannelID, embed)
}
