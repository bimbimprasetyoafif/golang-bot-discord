package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	r "golang-bot/response"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

const (
	GET     command = 1
	DEL     command = 2
	UPDATE  command = 3
	POST    command = 4
	Help    command = 5
	Nothing command = 0
)

type (
	command        int
	messageContent struct {
		Prefix string
		Args   string
	}
)

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Cannot read :" + err.Error())
	}

}

func main() {
	dg, err := discordgo.New(fmt.Sprintf("Bot %s", viper.GetString("token")))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageHandler)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	content := splitPrefixArgs(m.Content)
	switch checkCommandType(content.Prefix) {
	case GET:
		s.ChannelMessageSend(m.ChannelID, r.GetAllBook())
	case POST:
		s.ChannelMessageSend(m.ChannelID, r.CreateBook(content.Args))
	case UPDATE:
		s.ChannelMessageSend(m.ChannelID, r.UpdateBook(content.Args))
	case DEL:
		s.ChannelMessageSend(m.ChannelID, r.DeleteBook(content.Args))
	case Help:
		s.ChannelMessageSend(m.ChannelID, r.Help())
	}

}

func splitPrefixArgs(message string) messageContent {
	if !strings.Contains(message, "!") {
		return messageContent{
			Prefix: "",
			Args:   message,
		}
	}
	msgSplit := strings.Split(message, " ")
	return messageContent{
		Prefix: msgSplit[0],
		Args:   strings.Join(msgSplit[1:], " "),
	}
}

func checkCommandType(c string) command {
	switch c {
	case "!tambah":
		return POST
	case "!list":
		return GET
	case "!update":
		return UPDATE
	case "!hapus":
		return DEL
	case "!help":
		return Help
	default:
		return Nothing
	}
}
