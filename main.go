package main

import (
	"context"
	"os"
	"strings"

	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

)


func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
    if err != nil {
        log.Panic(err)
    }
	ctx := context.Background()
    bot.Debug = true

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil { // ignore any non-Message updates
            continue
        }

        if !update.Message.IsCommand() { // ignore any non-command Messages
            continue
        }

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		id := update.SentFrom().ID
		log.Print("Update from: ", id)
        switch update.Message.Command() {
        case "set":
			words := strings.Fields(update.Message.CommandArguments())
			if (len(words) < 2) {
            	msg.Text = "Not enough arguments"
			} else {
				login, password := words[0], words[1]
				log.Print("Setting: " , id, " ", login, " ", password)
				err := SetKey(ctx, Record{
					ID: id,
					Data: Data {
						Login: login,
						Password: password,
					},
				})
				if err != nil {
					msg.Text = "Error setting password"
					log.Print("Error: ", err)
				} else {
					msg.Text = "Success"
				}

			}
        case "get":
			words := strings.Fields(update.Message.CommandArguments())
			if (len(words) < 1) {
            	msg.Text = "Not enough arguments"
			} else {
				login := words[0]
				log.Print("Getting: " , login)
				val, err := GetKey(ctx, Record{
					ID: id,
					Data: Data {
						Login: login,
						Password: "",
					},
				})
				if err != nil {
					log.Print("Error: ", err)
					msg.Text = "Not found"
				} else {
					msg.Text = val
				}
			}
        case "del":
			words := strings.Fields(update.Message.CommandArguments())
			if (len(words) < 1) {
            	msg.Text = "Not enough arguments"
			} else {
				login := words[0]
				log.Print("Deleting" , login)
				err := DeleteKey(ctx, Record{
					ID: id,
					Data: Data {
						Login: login,
						Password: "",
					},
				})
				if err != nil {
					msg.Text = "Not found"
				} else {
					msg.Text = "Success"
				}
			}
        default:
            msg.Text = "I don't know that command"
        }

        if _, err := bot.Send(msg); err != nil {
            log.Panic(err)
        }
    }
}
