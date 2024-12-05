package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/types"
	"strconv"
)

func getMailboxCount(ctx context.Context, ms msgServer) (int, error) {
	it, err := ms.k.Mailboxes.Iterate(ctx, nil)
	if err != nil {
		return 0, err
	}

	mKeys, err := it.Keys()
	if err != nil {
		return 0, err
	}

	return len(mKeys), nil
}

// TODO: Add creator or domain ID
func generateMailboxID(count int) string {
	mailboxCount := []byte(strconv.Itoa(count))
	id := sha256.Sum256(mailboxCount)
	return "0x" + hex.EncodeToString(id[:])
}

func (ms msgServer) CreateMailbox(ctx context.Context, req *types.MsgCreateMailbox) (*types.MsgCreateMailboxResponse, error) {
	mailboxCount, err := getMailboxCount(ctx, ms)
	if err != nil {
		return nil, err
	}

	prefixedId := generateMailboxID(mailboxCount)

	newMailbox := types.Mailbox{
		Id:              prefixedId,
		Ism:             req.Ism,
		MessageSent:     0,
		MessageReceived: 0,
		Creator:         req.Creator,
	}

	if err = ms.k.Mailboxes.Set(ctx, prefixedId, newMailbox); err != nil {
		return nil, err
	}

	return &types.MsgCreateMailboxResponse{}, nil
}
