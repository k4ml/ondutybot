package main

import (
	"os"
	"log"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func RemoveIndex(s []tgbotapi.User, index int) []tgbotapi.User {
	return append(s[:index], s[index+1:]...)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	pendingNewUsers := make([]tgbotapi.User, 0)
	question := "Siapa PM kita yang ke 4 dan 7?"
	choices := "tar/dr.m, pak lah & dsn, dr.m & tun m"
	answer := "dr. m & tun m"
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		fmt.Println("New members", update.Message.NewChatMembers)
		if update.Message.NewChatMembers != nil {
			for _, newUser := range *update.Message.NewChatMembers {
				pendingNewUsers = append(pendingNewUsers, newUser)
				message := "Welcome %s. Sila jawab dgn Reply: %s?[%s]"
				message = fmt.Sprintf(message, newUser.FirstName, question, choices)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, message))
			}
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		fmt.Println("Pending users", pendingNewUsers)
		for index, pendingUser := range pendingNewUsers {
			if *update.Message.From == pendingUser {
				if update.Message.Text != answer {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong answer. You'll be removed."))
					chatMemberConfig := tgbotapi.ChatMemberConfig{ChatID: update.Message.Chat.ID, UserID: update.Message.From.ID}
					kickChatMemberConfig := tgbotapi.KickChatMemberConfig{ChatMemberConfig:chatMemberConfig}
					fmt.Println(kickChatMemberConfig)
					bot.KickChatMember(kickChatMemberConfig)
					pendingNewUsers = RemoveIndex(pendingNewUsers, index)
				} else {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Correct. You're verified"))
				}
			}
		}

		//fmt.Println(tgbotapi.ChatConfig(*update.Message.Chat))

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		//bot.Send(msg)
	}
}
