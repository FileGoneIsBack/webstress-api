package handler

import (
	"log"
	"bot/core/commands"
	send "bot/models/bot"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type CommandHandler func(bot *tgbotapi.BotAPI, message *tgbotapi.Message)
type ButtonHandler func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, callbackData string, userID int64)

var buttonHandlers = map[string]ButtonHandler{
	//"addbal": cmds.PaymentGateway, 
	"info":	cmds.Info,
	//"status":	cmds.Status,
	"auth":	cmds.Auth,
}

var commandHandlers = map[string]CommandHandler{
	"start": cmds.HandleHelpCommand,
}

func HandleUpdates(bot *tgbotapi.BotAPI) {
	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})
	if err != nil {
		log.Fatalf("Error getting updates: %v\n\n", err)
	}

	for update := range updates {
		if update.Message != nil {
			if command := update.Message.Command(); command != "" {
				if handler, found := commandHandlers[command]; found {
					handler(bot, update.Message)
				} else {
					log.Printf("No handler found for command: /%s", command)
					send.SendBounce(bot, update.Message, "Sorry, I don't understand that command.")
				}
			} else {
				//HandleUserInput(bot, update.Message)
			}
		} else if update.CallbackQuery != nil {
			log.Printf("Received callback data: %s\n\n", update.CallbackQuery.Data)

			userID := int64(update.CallbackQuery.From.ID)
			if handler, found := buttonHandlers[update.CallbackQuery.Data]; found {
				handler(bot, update.CallbackQuery.Message, update.CallbackQuery.Data, userID) 
			} else {
				send.SendBounce(bot, update.CallbackQuery.Message, "Unknown option")
			}

			_, err := bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
				CallbackQueryID: update.CallbackQuery.ID,
			})
			if err != nil {
				log.Printf("Error acknowledging callback: %v\n\n", err)
			}
		}
	}
}

/*
func HandleUserInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
    log.Printf("Handling user input: %s\n", message.Text)
	
    if state, exists := userStates[message.From.ID]; exists && state == WaitingForBalance {
        amountStr := message.Text
        amount, err := strconv.Atoi(amountStr)

        log.Printf("Parsed amount: %d\n", amount)
        if err != nil || amount <= 0 {
            reply := "Invalid amount. Please enter a valid number greater than 0."
            send.SendBounce(bot, message, reply)
            return
        }

        log.Printf("Sending invoice for amount: %d\n", amount)
        SendInvoice(bot, message, amount)

        userStates[message.From.ID] = DefaultState
    }
}
*/ 
