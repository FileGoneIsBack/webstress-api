package cmds

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"bot/models/sessions"
	"bot/database"
	send "bot/models/bot"
	"log"
)

func HandleHelpCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
    reply := "Welcome to the bot! Here are the available commands:\n"
    

    user, err := database.Container.GetUser(message.From.UserName)
    if err != nil {
        log.Printf("Error retrieving user: %v", err)
        reply = "An error occurred while fetching your data. Please try again later."
        send.SendBounce(bot, message, reply)
        return
    }
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

    // Check if the Telegram ID doesn't match the one in the database
    if fmt.Sprintf("%s", user.Telegram) != fmt.Sprintf("%d", message.From.ID) {
        // Create and store the session for the user
        session := &sessions.Session{
            ID:         message.From.ID,
            Username:   user.Username,
            Telegram:   user.Telegram,
            Balance:    user.Balance,
            Membership: user.Membership, 
        }

        sessions.SetSession(message.From.ID, session)
        reply = fmt.Sprintf("Hello, %s! Your account has been successfully linked.", user.Username)
        
        inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
            tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("Add bal", "addbal"),
                tgbotapi.NewInlineKeyboardButtonData("Account Info", "info"),
            ),
            tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("logs", "test3"),
                tgbotapi.NewInlineKeyboardButtonData("status", "test4"),
            ),
        )

        send.SendBounceWithInlineButtons(bot, message, reply, inlineKeyboard)
    } else {
        reply = "Your Telegram username is already linked with your account."
        send.SendBounce(bot, message, reply)
    }
}
