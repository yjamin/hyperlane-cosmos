package types

import (
	"context"
	"github.com/KYVENetwork/hyperlane-cosmos/util"
)

// combine multiple mailbox hooks, all hook functions are run in array sequence
var _ MailboxHooks = MultiMailboxHooks{}

type MultiMailboxHooks []MailboxHooks

func NewMultiMailboxHooks(hooks ...MailboxHooks) MultiMailboxHooks {
	return hooks
}

func (h MultiMailboxHooks) Handle(ctx context.Context, origin uint32, sender util.HexAddress, body []byte) error {
	for i := range h {
		if err := h[i].Handle(ctx, origin, sender, body); err != nil {
			return err
		}
	}

	return nil
}
