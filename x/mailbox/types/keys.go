package types

import "cosmossdk.io/collections"

const ModuleName = "mailbox"

var (
	ParamsKey            = collections.NewPrefix(0)
	MailboxesKey         = collections.NewPrefix(1)
	MailboxesSequenceKey = collections.NewPrefix(2)
	MessagesKey          = collections.NewPrefix(3)
)

var (
	// TODO: Set this dynamically.
	Domain = 1
)

var Version uint8 = 1
