package keeper

import (
	"context"
	"github.com/KYVENetwork/hyperlane-cosmos/util"
	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/types"
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

func (ms msgServer) CreateMailbox(ctx context.Context, req *types.MsgCreateMailbox) (*types.MsgCreateMailboxResponse, error) {
	mailboxCount, err := getMailboxCount(ctx, ms)
	if err != nil {
		return nil, err
	}

	prefixedId := util.CreateHexAddress(types.ModuleName, int64(mailboxCount))

	newMailbox := types.Mailbox{
		Id:              prefixedId.String(),
		Ism:             req.Ism,
		MessageSent:     0,
		MessageReceived: 0,
		Creator:         req.Creator,
	}

	if err = ms.k.Mailboxes.Set(ctx, prefixedId.Bytes(), newMailbox); err != nil {
		return nil, err
	}

	return &types.MsgCreateMailboxResponse{}, nil
}
