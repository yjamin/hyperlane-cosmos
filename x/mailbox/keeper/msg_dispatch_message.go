package keeper

import (
	"context"
	"github.com/KYVENetwork/hyperlane-cosmos/util"
	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (ms msgServer) DispatchMessage(ctx context.Context, req *types.MsgDispatchMessage) (*types.MsgDispatchMessageResponse, error) {

	goCtx := sdk.UnwrapSDKContext(ctx)

	bodyBytes, err := hexutil.Decode(req.Body)
	if err != nil {
		return nil, err
	}

	mailBoxId, err := util.DecodeHexAddress(req.MailboxId)
	if err != nil {
		return nil, err
	}

	sender, err := util.DecodeHexAddress(req.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := util.DecodeHexAddress(req.Recipient)
	if err != nil {
		return nil, err
	}

	msgId, err := ms.k.DispatchMessage(goCtx, mailBoxId, req.Destination, recipient, sender, bodyBytes)
	if err != nil {
		return nil, err
	}

	return &types.MsgDispatchMessageResponse{
		MessageId: msgId.String(),
	}, nil
}
