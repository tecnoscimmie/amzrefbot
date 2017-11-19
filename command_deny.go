package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Deny is the implementation of the "/deny" command
func (r *Refs) Deny(username string) {
	pUser, err := r.GetPendingUserByUsername(username)
	if err != nil {
		message := tgbotapi.NewMessage(userAdministrator, "User @"+pUser.AssociatedUser+" not found.")
		_, err = botInstance.Send(message)
		checkError(err)
		return
	}

	// compose the "sorry" message
	message := tgbotapi.NewMessage(pUser.ChatID, DENIED)
	_, err = botInstance.Send(message)
	checkError(err)

	// delete the pending user
	r.RemovePendingUser(username)
}
