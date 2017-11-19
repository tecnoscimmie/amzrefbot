package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Allow is the implementation of the "/allow" command
func (r *Refs) Allow(username string) {
	pUser, err := r.GetPendingUserByUsername(username)
	if err != nil {
		message := tgbotapi.NewMessage(userAdministrator, "User @"+username+" not found.")
		_, err = botInstance.Send(message)
		checkError(err)
		return
	}

	// compose the "welcome" message
	message := tgbotapi.NewMessage(pUser.ChatID, Accepted)
	_, err = botInstance.Send(message)
	checkError(err)

	// add a new user and delete it from the pending list
	r.SaveNewRefCode(pUser)
	r.RemovePendingUser(username)

}
