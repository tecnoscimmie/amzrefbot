package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

var (
	// Save refcodes in the local directory
	refSavePath = "./refcodes.json"
)

// Refs is a struct containing referral codes associated on a per-user basis.
// User will be able to add their data by sending messages to this bot.
type Refs struct {
	RefSavePath   string
	ReferralCodes []ReferralCode `json:"referral_codes"`
	// PendingUsers is an array of pending ReferralCode
	PendingUsers []ReferralCode
}

// NewRefs creates a Refs struct populated with the default save path.
// It tries to load from the default path an alread-existing referral code database.
func NewRefs() Refs {
	r := Refs{
		RefSavePath: refSavePath,
	}
	checkError(r.loadFromFile())
	return r
}

// NewRefsFromPath creates a Refs struct populated with the content of the file pointed by path.defa
// It tries to load from the path an alread-existing referral code database.
func NewRefsFromPath(path string) Refs {
	r := Refs{
		RefSavePath: path,
	}
	checkError(r.loadFromFile())
	return r
}

// ReferralCode is a single referral code, with data about its proprietary, the AssociatedUser
// and the Code itself.
type ReferralCode struct {
	AssociatedUser string `json:"associated_user"`
	Code           string `json:"code"`
	ChatID         int64  `json:"chat_id"`
}

// SaveNewRefCode saves the newly-added referral code to the struct, and saves it to disk.
func (r *Refs) SaveNewRefCode(rf ReferralCode) {
	r.ReferralCodes = append(r.ReferralCodes, rf)
	// saveToFile in a goroutine, non-blocking call
	go func() {
		err := r.saveToFile()
		if err != nil {
			log.Println("error during database write:", err)
		}
	}()
}

// saveToFile saves its struct to file.
func (r *Refs) saveToFile() (err error) {
	codesByte, err := json.Marshal(r)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(r.RefSavePath, codesByte, 0644)
	return
}

// loadFromFile loads its struct from file.
func (r *Refs) loadFromFile() (err error) {
	fileContent, err := ioutil.ReadFile(r.RefSavePath)
	if err != nil {
		return
	}

	// unmarshal json into the file
	var newR Refs
	err = json.Unmarshal(fileContent, &newR)
	if err != nil {
		return
	}

	r.ReferralCodes = newR.ReferralCodes
	r.RefSavePath = newR.RefSavePath

	return
}

// PrettyPrintList returns a pretty-printed string containing all the current pending requests
func (r Refs) PrettyPrintList() (result string) {
	if len(r.PendingUsers) <= 0 {
		result = fmt.Sprintf("There are no requests currently awaiting approval!\nGood job :)\n")
		return
	}

	result = fmt.Sprintf("There are currently %d requests awaiting approval:\n\n", len(r.PendingUsers))
	for index, item := range r.PendingUsers {
		result = result + fmt.Sprintf("%d) @%s\n", index+1, item.AssociatedUser)
	}

	result = result + fmt.Sprintf("\nUse the /deny {username} or /allow {username} commands to either deny or allow the addition to the database!")
	return
}

// GetPendingUserByUsername returns a pending user by its username.
func (r Refs) GetPendingUserByUsername(username string) (pendingUser ReferralCode, err error) {
	for _, puser := range r.PendingUsers {
		if strings.ToUpper(puser.AssociatedUser) == strings.ToUpper(username) {
			pendingUser = puser
			return
		}
	}

	err = errors.New(fmt.Sprintln(username, "not found"))
	return
}

// GetPendingUserByID returns a pending user by its id.
func (r Refs) GetPendingUserByID(id int64) (pendingUser ReferralCode, err error) {
	for _, puser := range r.PendingUsers {
		if puser.ChatID == id {
			pendingUser = puser
			return
		}
	}

	err = errors.New(fmt.Sprintln(id, "not found"))
	return
}

// GetUserByUsername returns a user by its username.
func (r Refs) GetUserByUsername(username string) (user ReferralCode, err error) {
	for _, ruser := range r.ReferralCodes {
		if strings.ToUpper(ruser.AssociatedUser) == strings.ToUpper(username) {
			user = ruser
			return
		}
	}

	err = errors.New(fmt.Sprintln(username, "not found"))
	return
}

// GetUserByID returns a pending user by its id.
func (r Refs) GetUserByID(id int64) (user ReferralCode, err error) {
	for _, ruser := range r.PendingUsers {
		if ruser.ChatID == id {
			user = ruser
			return
		}
	}

	err = errors.New(fmt.Sprintln(id, "not found"))
	return
}

// AlreadyGotUser checks if the user is already either in the pending list or on the reflist
func (r Refs) AlreadyGotUser(username string, id int64) bool {
	_, err := r.GetUserByUsername(username)
	_, errID := r.GetUserByID(id)
	_, errPending := r.GetPendingUserByUsername(username)
	_, errPendingID := r.GetPendingUserByID(id)

	if err != nil && errPending != nil {
		return false
	}

	if errID != nil && errPendingID != nil {
		return false
	}

	return true
}

// RemovePendingUser remove a pending request from the session.
func (r *Refs) RemovePendingUser(username string) {
	for index, user := range r.PendingUsers {
		if strings.ToUpper(user.AssociatedUser) == strings.ToUpper(username) {
			r.PendingUsers[len(r.PendingUsers)-1], r.PendingUsers[index] = r.PendingUsers[index], r.PendingUsers[len(r.PendingUsers)-1]
			r.PendingUsers = r.PendingUsers[:len(r.PendingUsers)-1]
		}
	}

}
