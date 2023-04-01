package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	TelegramToken string
	bot           *tgbotapi.BotAPI
)

func setTask(mins int, out chan int) {
	defer close(out)

	for i := 0; i < mins; i++ {
		out <- i
		time.Sleep(time.Second)
	}
}

func handleTask(update tgbotapi.Update) {
	words := strings.Split(update.Message.Text, " ")
	if len(words) < 2 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide number of seconds along with the command")
		bot.Send(msg)
		return
	}

	mins, _ := strconv.Atoi(words[1])
	messageText := fmt.Sprintf("You have set a timer for %d seconds", mins)
	message := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
	bot.Send(message)

	out := make(chan int)
	go setTask(mins, out)
	for i := range out {
		messageTextTimer := fmt.Sprintf("%s\nYou have %d seconds left 🕥", message.Text, mins-i)
		editedMessage := tgbotapi.NewEditMessageText(update.Message.Chat.ID, update.Message.MessageID+1, messageTextTimer)
		bot.Send(editedMessage)
	}

	messageText = messageText + "\nYour task has finished 🎉"
	editedMessage := tgbotapi.NewEditMessageText(update.Message.Chat.ID, update.Message.MessageID+1, messageText)
	bot.Send(editedMessage)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	TelegramToken = os.Getenv("TELEGRAM_TOKEN")
}

func main() {
	var err error
	bot, err = tgbotapi.NewBotAPI(TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "task":
					go handleTask(update)
				}
			}
		}
	}
}
