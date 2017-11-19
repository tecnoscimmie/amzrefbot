package main

const (
	// AddCodeCommand - command used to add a code to the database.
	// Any user can use this command, but an administrator must approve or deny the request.
	AddCodeCommand = "/addcode"

	// DenyCommand - command used to deny the addidion of a code to the database.
	// Only an administrator can use this command.
	DenyCommand = "/deny"

	// AllowCommand - command used to allow the addition og a code to the database.
	// Only an administrator can use this command.
	AllowCommand = "/allow"

	// ListPendingCommand - command used to list all the pending requests currently in memory.
	// Only an administrator can use this command
	ListPendingCommand = "/listpending"
)
