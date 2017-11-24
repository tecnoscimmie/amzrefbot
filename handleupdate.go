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
		message := strings.Split(update.Message.Text, " ")
		debugPrint(fmt.Sprintln("parsed message array:", message))
		debugPrint("command -> " + message[0])
		switch message[0] { // the first element of the array is usually the command
		case ListPendingCommand:
			debugPrint("called handleListPending()")
			handleListPending(int64(update.Message.From.ID))
		case DenyCommand:
			adminErr := checkIfAdministrator(update.Message.From.ID)
			paramsErr := checkIfHasParameter(message, update.Message.From.ID, "deny")
			if adminErr != nil || paramsErr != nil {
				return
			}
			debugPrint("called handleDeny()")
			handleDeny(message[1])
		case AllowCommand:
			adminErr := checkIfAdministrator(update.Message.From.ID)
			paramsErr := checkIfHasParameter(message, update.Message.From.ID, "allow")
			if adminErr != nil || paramsErr != nil {
				return
			}
			debugPrint("called handleAllow()")
			handleAllow(message[1])
		case AddCodeCommand:
			paramsErr := checkIfHasParameter(message, update.Message.From.ID, "addcode")
			if paramsErr != nil {
				return
			}
			debugPrint("called handleAddCode()")
			handleAddCode(update.Message.From.UserName, message[1], int64(update.Message.From.ID))
		default:
			debugPrint("called handleDefault()")
			handleDefault(int64(update.Message.From.ID))
		}
	} else if update.CallbackQuery != nil {
		log.Println("got callback query with message:", update.CallbackQuery.Data)
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
