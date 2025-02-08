package cmds

import (
	//"bot/database"
	send "bot/models/bot"
	"bot/models/sessions"
	"log"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func Info(bot *tgbotapi.BotAPI, message *tgbotapi.Message, callbackData string, userID int) {
	user, _ := sessions.GetSession(userID)

	infoMessage := "User Info:\n"
	infoMessage += "ID: " + strconv.Itoa(user.ID) + "\n"
	infoMessage += "User: " + user.Username + "\n"
	infoMessage += "TGUser: " + user.Telegram + "\n"
	infoMessage += "Bal: " + strconv.Itoa(user.Balance) + "\n"
	infoMessage += "Membership: " + user.Membership + "\n"


	send.SendBounce(bot, message, infoMessage)
	log.Printf("User %d info sent: %s", userID, infoMessage)
}
