package keeper

import (
	"context"
	"errors"
	"github.com/KYVENetwork/hyperlane-cosmos/util"
	"github.com/KYVENetwork/hyperlane-cosmos/x/warp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	k Keeper
}

func (ms msgServer) CreateCollateralToken(ctx context.Context, msg *types.MsgCreateCollateralToken) (*types.MsgCreateCollateralTokenResponse, error) {

	next, err := ms.k.Sequence.Next(ctx)
	if err != nil {
		return nil, err
	}

	mailboxId, err := util.DecodeHexAddress(msg.OriginMailbox)
	if err != nil {
		return nil, err
	}

	has, err := ms.k.mailboxKeeper.Mailboxes.Has(ctx, mailboxId.Bytes())
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("mailbox not found")
	}

	receiverMailbox, err := util.DecodeHexAddress(msg.ReceiverMailbox)
	if err != nil {
		return nil, err
	}

	receiverContract, err := util.DecodeHexAddress(msg.ReceiverContract)
	if err != nil {
		return nil, err
	}

	tokenId := util.CreateHexAddress(types.ModuleName, int64(next))

	newToken := types.HypToken{
		Id:               tokenId.Bytes(),
		Creator:          msg.Creator,
		TokenType:        types.HYP_TOKEN_COLLATERAL,
		OriginMailbox:    mailboxId.Bytes(),
		OriginDenom:      msg.OriginDenom,
		ReceiverDomain:   msg.ReceiverDomain,
		ReceiverMailbox:  receiverMailbox.Bytes(),
		ReceiverContract: receiverContract.Bytes(),
	}

	if err = ms.k.HypTokens.Set(ctx, tokenId.Bytes(), newToken); err != nil {
		return nil, err
	}
	return &types.MsgCreateCollateralTokenResponse{}, nil
}

func (ms msgServer) RemoteTransfer(ctx context.Context, msg *types.MsgRemoteTransfer) (*types.MsgRemoteTransferResponse, error) {
	goCtx := sdk.UnwrapSDKContext(ctx)

	tokenId, err := util.DecodeHexAddress(msg.TokenId)
	if err != nil {
		return nil, err
	}

	token, err := ms.k.HypTokens.Get(ctx, tokenId.Bytes())
	if err != nil {
		return nil, err
	}

	senderAcc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	err = ms.k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAcc, types.ModuleName, sdk.NewCoins(sdk.NewCoin(token.OriginDenom, msg.Amount)))
	if err != nil {
		return nil, err
	}

	recipient, err := util.DecodeEthHex(msg.Recipient)
	if err != nil {
		return nil, err
	}

	warpPayload, err := types.NewWarpPayload(recipient, *msg.Amount.BigInt())
	if err != nil {
		return nil, err
	}

	// Token destinationDomain, recipientAddress
	dispatchMsg, err := ms.k.mailboxKeeper.DispatchMessage(
		goCtx,
		util.HexAddress(token.OriginMailbox),
		token.ReceiverDomain,
		util.HexAddress(token.ReceiverContract),
		tokenId,
		warpPayload.Bytes(),
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgRemoteTransferResponse{
		MessageId: dispatchMsg.String(),
	}, nil
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}
