package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// debug log.Println's str passed in input only if the global "debug" variable is true
func debugPrint(str string) {
	if debug {
		log.Println(str)
	}
}

// checkError just checks and logs an error to stderr if any
func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func checkIfAdministrator(chatid int) error {
	chatid64 := int64(chatid)
	if userAdministrator != chatid64 {
		message := tgbotapi.NewMessage(chatid64, "You're not authorized to use this command.")
		_, err := botInstance.Send(message)
		checkError(err)
		return errors.New("user not authorized")
	}
	return nil
}

func checkIfHasParameter(message []string, chatid int, funcname string) error {
	chatid64 := int64(chatid)
	if len(message) < 2 {
		var strmsg string
		switch funcname {
		case "/deny", "/allow":
			strmsg = fmt.Sprintf("Missing username!\nSyntax: `/%s username`\n", funcname)
		case "/addcode":
			strmsg = fmt.Sprintf("Missing username!\nSyntax: `/%s refcode`\n", funcname)
		}
		tgMessage := tgbotapi.NewMessage(chatid64, strmsg)
		tgMessage.ParseMode = tgbotapi.ModeMarkdown
		_, err := botInstance.Send(tgMessage)
		checkError(err)
		return errors.New("missing parameter")
	}
	return nil
}
