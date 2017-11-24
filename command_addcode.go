package main

import "github.com/go-telegram-bot-api/telegram-bot-api"

// This file describes and implements the flow of execution of the "/addcode" command.
//
// When a non-administrator user sends a "/addcode" command:
// 	1. the system will parse the message, and assumes all the string part after the command itself
//	   as the code to be added. The username is the message sender.
//     A new ReferralCode struct will be created with the data, and added to the Pending struct.
//	2. the system will send a message to the user, saying that its add request needs to be approved
//	   by an administrator. The system will then send a message to all the administrators, reporting a new
//     user asked to be added to the referral code database.
//
// The administrator can then use the "/allow" or "/deny" command to allow or deny an user.

// AddCode adds a refcode to the pending user list.
func (r *Refs) AddCode(username string, refcode string, chatid int64) {
	// check if the user is already in the database
	if r.AlreadyGotUser(username, chatid) {
		// send message to the user
		message := tgbotapi.NewMessage(chatid, AlreadySent)
		_, err := botInstance.Send(message)
		checkError(err)
		return
	}

	r.PendingUsers = append(r.PendingUsers, ReferralCode{AssociatedUser: username, Code: refcode, ChatID: chatid})

	// send message to the user
	t := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Allow @"+username, "/allow "+username),
		tgbotapi.NewInlineKeyboardButtonData("Deny @"+username, "/deny "+username),
	)
	message := tgbotapi.NewMessage(chatid, ApproveOrDeny)
	message.ReplyMarkup = t
	_, err := botInstance.Send(message)
	checkError(err)

	// send message to the administrators
	message = tgbotapi.NewMessage(userAdministrator, GetAdminAddNotify(username))
	_, err = botInstance.Send(message)
	checkError(err)
}
