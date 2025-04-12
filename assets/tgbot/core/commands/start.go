package cmds

import (
	"bot/database"
	send "bot/models/bot"
	"bot/models/sessions"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandleHelpCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
    reply := "Welcome to the bot! Here are the available commands:\n"
    log.Printf("[DEBUG] User joined - TG ID: %d", int64(message.From.ID))
    user, _ := database.Container.GetUser(int64(message.From.ID))
    if user == nil {
        reply = "No user found with your Telegram username. Please register first."
        
        inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
            tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("Auth", "auth"), 
            ),
        )
        send.SendBounceWithInlineButtons(bot, message, reply, inlineKeyboard)
        return
    }
    session := &sessions.Session{
        ID:         message.From.ID,
        Username:   user.Username,
        Telegram:   user.Tele,
        Balance:    user.Balance,
        Membership: user.Membership, 
    }
    sessions.SetSession(message.From.ID, session)

    reply += "\nYou're logged in as: " + user.Username
    inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Info", "info"), 
        ),
    )
    send.SendBounceWithInlineButtons(bot, message, reply, inlineKeyboard)
}