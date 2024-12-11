package keeper

import (
	"context"
	"fmt"
	"github.com/KYVENetwork/hyperlane-cosmos/util"
	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/types"
)

func (ms msgServer) CreateMailbox(ctx context.Context, req *types.MsgCreateMailbox) (*types.MsgCreateMailboxResponse, error) {
	mailboxCount, err := ms.k.MailboxesSequence.Next(ctx)
	if err != nil {
		return nil, err
	}

	ismExists, err := ms.k.ismKeeper.IsmIdExists(ctx, req.Ism)
	if err != nil {
		return nil, err
	}
	if !ismExists {
		return nil, fmt.Errorf("ISM %s doesn't exist", req.Ism)
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
