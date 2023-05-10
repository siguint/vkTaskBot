package main

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

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

	ticker := time.NewTicker(15 * time.Second)
	mu := sync.Mutex{}
	messagesToDelete := make([]tgbotapi.DeleteMessageConfig, 0)

	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				mu.Lock()
				for _, x := range messagesToDelete {
					_, err := bot.Request(x)

					if err != nil {
						log.Print(err)
					}
				}
				messagesToDelete = nil
				mu.Unlock()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		id := update.SentFrom().ID
		mu.Lock()
		messagesToDelete = append(messagesToDelete, tgbotapi.DeleteMessageConfig{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.MessageID,
		})
		mu.Unlock()
		log.Print("Update from: ", id)
		switch update.Message.Command() {
		case "set":
			words := strings.Fields(update.Message.CommandArguments())
			if len(words) < 3 {
				msg.Text = "Not enough arguments"
			} else {
				service, login, password := words[0], words[1], words[2]
				err := SetKey(ctx, Record{
					ID:      id,
					Service: service,
					Data: Data{
						Login:    login,
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
			if len(words) < 1 {
				msg.Text = "Not enough arguments"
			} else {
				service := words[0]
				log.Print("Getting: ", service)
				val, err := GetKey(ctx, Record{
					ID:      id,
					Service: service,
					Data: Data{
						Login:    "",
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
			if len(words) < 1 {
				msg.Text = "Not enough arguments"
			} else {
				service := words[0]
				log.Print("Deleting", service)
				err := DeleteKey(ctx, Record{
					ID:      id,
					Service: service,
					Data: Data{
						Login:    "",
						Password: "",
					},
				})
				if err != nil {
					msg.Text = "Not found"
				} else {
					msg.Text = "Success"
				}
			}
		case "help":
			msg.Text = `
				Commands:
					/set login password - adds login and password to service
					/get login - get login and password by service name
					/del login - delete login and password by service name
			`
		default:
			msg.Text = "I don't know that command"
		}

		if msg, err := bot.Send(msg); err != nil {
			log.Panic(err)
		} else {
			mu.Lock()
			messagesToDelete = append(messagesToDelete, tgbotapi.DeleteMessageConfig{
				ChatID:    msg.Chat.ID,
				MessageID: msg.MessageID,
			})
			mu.Unlock()
		}
	}
	close(quit)
}
