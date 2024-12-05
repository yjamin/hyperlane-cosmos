package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/types"
)

func (ms msgServer) ProcessMessage(ctx context.Context, req *types.MsgProcessMessage) (*types.MsgProcessMessageResponse, error) {
	rawMessage := hexutil.MustDecode(req.Message)

	// TODO: Use Metadata for ISM.
	_ = hexutil.MustDecode(req.Metadata)

	// TODO: Check if destination domain is current domain
	version := Version(rawMessage)
	if version != types.Version {
		return nil, fmt.Errorf("mailbox: bad version %v", version)
	}

	rawRecipient := Recipient(rawMessage)
	recipient := hexutil.Encode(rawRecipient)

	messageId := Id(rawMessage)
	key := collections.Join(recipient, messageId)

	received, err := ms.k.Messages.Has(ctx, key)
	if err != nil {
		return nil, err
	}
	if received {
		return nil, fmt.Errorf("already received messsage")
	}

	mailbox, err := ms.k.Mailboxes.Get(ctx, recipient)
	if err != nil {
		return nil, err
	}

	mailbox.MessageReceived++

	err = ms.k.Mailboxes.Set(ctx, mailbox.Id, mailbox)
	if err != nil {
		return nil, err
	}

	origin := Origin(rawMessage)
	rawSender := Sender(rawMessage)
	sender := hexutil.Encode(rawSender)
	body := Body(rawMessage)

	// TODO: Process Message.

	err = ms.k.Messages.Set(ctx, key)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// TODO: Add OriginMailboxId
	// TODO: Check if OriginMailboxId should be used as sender
	_ = sdkCtx.EventManager().EmitTypedEvent(&types.Process{
		OriginMailboxId: "",
		OriginDomain:    origin,
		Sender:          sender,
		Recipient:       recipient,
		MessageId:       hexutil.Encode(messageId),
		MessageBody:     hexutil.Encode(body),
	})

	return &types.MsgProcessMessageResponse{}, nil
}
