package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"encoding/binary"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/types"
)

func (ms msgServer) DispatchMessage(ctx context.Context, req *types.MsgDispatchMessage) (*types.MsgDispatchMessageResponse, error) {
	mailbox, err := ms.k.Mailboxes.Get(ctx, req.MailboxId)
	if errors.Is(err, collections.ErrNotFound) {
		// TODO: Handle
	}
	if err != nil {
		return nil, err
	}

	mailbox.MessageSent++

	err = ms.k.Mailboxes.Set(ctx, req.MailboxId, mailbox)
	if err != nil {
		return nil, err
	}

	var message []byte

	message = append(message, []byte{types.Version}...)

	nonce := make([]byte, 4)
	binary.BigEndian.PutUint32(nonce, mailbox.MessageSent)
	message = append(message, nonce...)

	domain := make([]byte, 4)
	binary.BigEndian.PutUint32(domain, uint32(types.Domain))
	message = append(message, domain...)

	// TODO: Check if this should be mailbox ID
	sender := sdk.MustAccAddressFromBech32(req.Sender).Bytes()
	for len(sender) < (DestinationOffset - SenderOffset) {
		padding := make([]byte, 1)
		sender = append(padding, sender...)
	}
	message = append(message, sender...)

	destinationDomain := make([]byte, 4)
	binary.BigEndian.PutUint32(domain, req.Destination)
	message = append(message, destinationDomain...)

	// TODO: Check if hex is correct encoding
	recipient := hexutil.MustDecode(req.Recipient)
	if err != nil {
		return nil, err
	}
	for len(recipient) < (BodyOffset - RecipientOffset) {
		padding := make([]byte, 1)
		recipient = append(padding, recipient...)
	}
	message = append(message, recipient...)

	// TODO: Check if hex is correct encoding
	// TODO: Verify if body matches expected format.
	body := hexutil.MustDecode(req.Body)
	if err != nil {
		return nil, err
	}
	message = append(message, body...)
	messageId := crypto.Keccak256(message)

	key := collections.Join(req.MailboxId, messageId)

	err = ms.k.Messages.Set(ctx, key)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_ = sdkCtx.EventManager().EmitTypedEvent(&types.Dispatch{
		DestinationDomain: req.Destination,
		// TODO: Use recipient
		RecipientAddress: req.Recipient,
		MessageBody:      hexutil.Encode(message),
		OriginDomain:     uint32(types.Domain),
		OriginMailbox:    req.MailboxId,
		// TODO: Verify type
		Sender: req.Sender,
	})

	return &types.MsgDispatchMessageResponse{
		MessageId: hexutil.Encode(messageId),
	}, nil
}
