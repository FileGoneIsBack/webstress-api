package cmds

import (
	send "bot/models/bot"
	"bot/database"
	"log"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func Auth(bot *tgbotapi.BotAPI, message *tgbotapi.Message, callbackData string, userID int64) {
	// Create an invite token for the user
	token, err := database.Container.CreateInvite(userID)
	if err != nil {
		send.SendBounce(bot, message, "Failed to create invite token.")
		log.Printf("Failed to create invite for user %d: %v", userID, err)
		return
	}

	signupURL := fmt.Sprintf("https://yourdomain.com/signup?token=%s", token)

	// Create button that links to the signup page
	button := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Sign Up", signupURL),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Here is your invite token:\n`%s`\n\nClick the button below to sign up:", token))
	msg.ReplyMarkup = button
	msg.ParseMode = "Markdown"

	_, _ = bot.Send(msg)

	log.Printf("User %d invite sent: %s", userID, token)
}