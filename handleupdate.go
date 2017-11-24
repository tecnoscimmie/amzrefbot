package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// handleUpdate handles updates sent from Telegram to us
func handleUpdate(update tgbotapi.Update) {
	if update.InlineQuery != nil {
		inputMsg := update.InlineQuery.Query
		debugPrint("received an inline query from " + update.InlineQuery.From.UserName)
		// inline query, check for amazon.* urls, get a random refcode and build the url
		val, user, err := refs.GenAffiliate(inputMsg)
		if !(err != nil || len(refs.ReferralCodes) < 0) {
			happyString := fmt.Sprintf("The lucky winner is is @%s!\nHere's the reflink: %s\n", user, val)
			okArticle := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID+"-ok", SuccessTitle, happyString)
			okArticle.Description = val

			inlineConf := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     0,
				Results:       []interface{}{okArticle},
			}

			if _, err := botInstance.AnswerInlineQuery(inlineConf); err != nil {
				log.Println(err)
			}
		}
	} else if update.Message != nil {
		debugPrint("received a private message from " + update.Message.From.UserName)
		// private message, handle various messages
		handleTelegramMessage(update.Message.Text, int64(update.Message.From.ID), update.Message.From.UserName)
	} else if update.CallbackQuery != nil {
		log.Println("got callback query with message:")
		handleTelegramMessage(update.CallbackQuery.Data, int64(update.CallbackQuery.From.ID), update.CallbackQuery.From.UserName)
	}
}

func handleTelegramMessage(messageBlob string, chatid int64, username string) {
	message := strings.Split(messageBlob, " ")
	debugPrint(fmt.Sprintln("parsed message array:", message))
	debugPrint("command -> " + message[0])
	switch message[0] { // the first element of the array is usually the command
	case ListPendingCommand:
		debugPrint("called handleListPending()")
		handleListPending(chatid)
	case DenyCommand:
		adminErr := checkIfAdministrator(int(chatid))
		paramsErr := checkIfHasParameter(message, int(chatid), "deny")
		if adminErr != nil || paramsErr != nil {
			return
		}
		debugPrint("called handleDeny()")
		handleDeny(message[1])
	case AllowCommand:
		adminErr := checkIfAdministrator(int(chatid))
		paramsErr := checkIfHasParameter(message, int(chatid), "allow")
		if adminErr != nil || paramsErr != nil {
			return
		}
		debugPrint("called handleAllow()")
		handleAllow(message[1])
	case AddCodeCommand:
		paramsErr := checkIfHasParameter(message, int(chatid), "addcode")
		if paramsErr != nil {
			return
		}
		debugPrint("called handleAddCode()")
		handleAddCode(username, message[1], chatid)
	default:
		debugPrint("called handleDefault()")
		handleDefault(chatid)
	}
}

func handleListPending(user int64) {
	listStr := refs.PrettyPrintList()
	var message tgbotapi.MessageConfig
	if user == userAdministrator {
		message = tgbotapi.NewMessage(user, listStr)
	} else {
		message = tgbotapi.NewMessage(user, "You're not authorized to use this command.")
	}

	_, err := botInstance.Send(message)
	checkError(err)
}

func handleDeny(user string) {
	refs.Deny(user)
}

func handleAllow(user string) {
	refs.Allow(user)
}

func handleAddCode(user string, refcode string, chatid int64) {
	refs.AddCode(user, refcode, chatid)
}

func handleDefault(user int64) {
	message := tgbotapi.NewMessage(user, "What? I didn't understood that command!")
	_, err := botInstance.Send(message)
	checkError(err)
}
