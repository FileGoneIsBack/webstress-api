package cmds

import (
	//"bot/database"
	send "bot/models/bot"
	"bot/models/sessions"
	"log"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func Info(bot *tgbotapi.BotAPI, message *tgbotapi.Message, callbackData string, userID int64) {
	user, _ := sessions.GetSession(int(userID))

	infoMessage := "User Info:\n"
	infoMessage += "ID: " + strconv.Itoa(user.ID) + "\n"
	infoMessage += "User: " + user.Username + "\n"
	infoMessage += "TGUser: " + strconv.FormatInt(user.Telegram, 10) + "\n"
	infoMessage += "Bal: " + strconv.Itoa(user.Balance) + "\n"
	infoMessage += "Membership: " + user.Membership + "\n"


	send.SendBounce(bot, message, infoMessage)
	log.Printf("User %d info sent: %s", userID, infoMessage)
}
