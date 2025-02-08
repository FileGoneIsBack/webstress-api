package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

// SendBounce 
func SendBounce(bot *tgbotapi.BotAPI, message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// SendBounceWithButtons
func SendBounceWithReply(bot *tgbotapi.BotAPI, message *tgbotapi.Message, text string, replyMarkup tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = replyMarkup

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// SendBounceWithInlineButtons
func SendBounceWithInlineButtons(bot *tgbotapi.BotAPI, message *tgbotapi.Message, text string, inlineMarkup tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = inlineMarkup

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}


