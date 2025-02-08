package cmds

import (
	"bot/models"
	send "bot/models/bot"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var userStates = make(map[int]string) 

const (
    WaitingForBalance = "waiting_for_balance"
    DefaultState      = "default"
)


func PaymentGateway(bot *tgbotapi.BotAPI, message *tgbotapi.Message, callbackData string, userID int) {
    reply := "Please enter the amount of balance you want to add (in USD):"
    send.SendBounce(bot, message, reply)
	
    userStates[userID] = WaitingForBalance
    log.Printf("User %d state set to: %s", userID, WaitingForBalance)
}

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


func SendInvoice(bot *tgbotapi.BotAPI, message *tgbotapi.Message, bal int) {
	title := "Finish Payment"
	description := "Purchase 1 credit for $1"
	currency := models.Config.Autobuy.Flat
	price := bal * 100 

	labeledPrice := []tgbotapi.LabeledPrice{
		{
			Label:  description,
			Amount: price,
		},
	}

	invoice := struct {
		ChatID         int64             `json:"chat_id"`
		Title          string            `json:"title"`
		Description    string            `json:"description"`
		Payload        string            `json:"payload"`
		ProviderToken  string            `json:"provider_token"`
		//StartParameter string            `json:"start_parameter"`
		Currency       string            `json:"currency"`
		Prices         []tgbotapi.LabeledPrice `json:"prices"`
	}{
		ChatID:         message.Chat.ID,
		Title:          title,
		Description:    description,
		Payload:        "unique_payload_here", 
		ProviderToken:  models.Config.Autobuy.Key, 
		Currency:       currency,
		Prices:         labeledPrice, 
	}

	data, err := json.Marshal(invoice)
	if err != nil {
		log.Printf("Error encoding invoice data: %v", err)
		return
	}

	url := "https://api.telegram.org/bot" + bot.Token + "/sendInvoice"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error sending invoice: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send invoice: %s", resp.Status)
	} else {
		log.Println("Invoice sent successfully!")
	}
}
