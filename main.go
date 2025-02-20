package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"discord-bot/commands"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Config struct {
	Token  string `json:token`
	Status string `json:status`
}

var token string
var config = loadConfig()

func main() {
	if token == "" {
		token = config.Token
	}

	if token == "" {
		fmt.Println("Bot tokeni bulunamadı, token girin: ")
		_, _ = fmt.Scanln(&token)
	}

	dc, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("Bot çalıştırılırken bir hata çıktı: ", err)
		return
	}

	dc.AddHandler(func(session *discordgo.Session, msg *discordgo.MessageCreate) {
		if session.State.User.ID == msg.Author.ID {
			return
		}

		if !strings.HasPrefix(msg.Content, "!") {
			return
		}

		text := strings.TrimPrefix(msg.Content, "!")
		
		parts := strings.Split(text, " ")
		name := parts[0]
		args := parts[1:]

		cmd, ok := commands.GetCommands()[name]
		if ok {
			cmd(session, msg.Message, args)
		} else {
			_, _ = session.ChannelMessageSend(msg.ChannelID, "Tanımsız komut: "+name)
		}
	})

	dc.Identify.Intents = discordgo.IntentsGuildMessages

	err = dc.Open()

	if err != nil {
		fmt.Println("Websocket açılırken hata oluştu, ", err)
		return
	}

	go func() {
		for {
			if config.Status != "" {
				dc.UpdateGameStatus(0, config.Status)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	fmt.Println("Bot çalışıyor.")
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-signalChannel

	_ = dc.Close()
}

func loadConfig() Config {
	cfg := Config{}
	_, err := os.Stat("config.json")
	if os.IsNotExist(err) {
		bytes, _ := json.MarshalIndent(cfg, "", "")
		_ = ioutil.WriteFile("config.json", bytes, 0644)
	} else {
		bytes, _ := ioutil.ReadFile("config.json")
		_ = json.Unmarshal(bytes, &cfg)
	}

	return cfg
}
