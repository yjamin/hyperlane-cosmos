package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"fmt"
	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated message.
	// typically, this should be the x/gov module account.
	authority string

	// state management
	Mailboxes collections.Map[string, types.Mailbox]
	Messages  collections.KeySet[collections.Pair[string, []byte]]
	Params    collections.Item[types.Params]
	Schema    collections.Schema
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		Mailboxes:    collections.NewMap(sb, types.MailboxesKey, "mailboxes", collections.StringKey, codec.CollValue[types.Mailbox](cdc)),
		Messages:     collections.NewKeySet(sb, types.MessagesKey, "messages", collections.PairKeyCodec(collections.StringKey, collections.BytesKey)),
		Params:       collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}
