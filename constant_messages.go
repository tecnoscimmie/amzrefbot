package main

const (
	// SuccessTitle is a generic success message for inline replies
	SuccessTitle = "Here's your reflink!"

	// PersonalTitle is a success message for personal reflinks
	PersonalTitle = "Personal reflink"

	// AlreadySent is an already sent request
	AlreadySent = "You already sent this request."

	// ApproveOrDeny is a message sent to a user when the registration is successful
	ApproveOrDeny = "An administrator will either approve or deny your request as soon as possible!"

	// HalfAdminAddNotify is a half admin add notify string
	HalfAdminAddNotify = "A new code add request has been received from @"

	// DENIED is.. DENIED
	DENIED = "The administrator decided that your code cannot be added to this bot.\nSorry!"

	// Accepted is happy feeling!
	Accepted = "Your request has been accepted!"
)

// GetAdminAddNotify builds a string with a notification for the admin for a new user
func GetAdminAddNotify(username string) string {
	return HalfAdminAddNotify + username + "!"
}
