package keeper

import (
	"context"
	"encoding/hex"
	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

func (ms msgServer) ProcessMessage(ctx context.Context, req *types.MsgProcessMessage) (*types.MsgProcessMessageResponse, error) {

	goCtx := sdk.UnwrapSDKContext(ctx)

	// Decode and parse message
	if strings.HasPrefix(req.Message, "0x") {
		req.Message = req.Message[2:]
	}
	messageBytes, err := hex.DecodeString(req.Message)
	if err != nil {
		return nil, err
	}

	// Decode and parse metadata
	if strings.HasPrefix(req.Metadata, "0x") {
		req.Metadata = req.Metadata[2:]
	}
	metadataBytes, err := hex.DecodeString(req.Metadata)
	if err != nil {
		return nil, err
	}

	if err = ms.k.ProcessMessage(goCtx, req.MailboxId, messageBytes, metadataBytes); err != nil {
		return nil, err
	}

	return &types.MsgProcessMessageResponse{}, nil
}
