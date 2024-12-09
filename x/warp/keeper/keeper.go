package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"encoding/hex"
	"fmt"
	"github.com/KYVENetwork/hyperlane-cosmos/util"
	mailboxkeeper "github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/keeper"
	"github.com/KYVENetwork/hyperlane-cosmos/x/warp/types"
	"github.com/cosmos/cosmos-sdk/codec"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated message.
	// typically, this should be the x/gov module account.
	authority string

	// state management
	HypTokens collections.Map[[]byte, types.HypToken]
	Params    collections.Item[types.Params]
	Schema    collections.Schema
	Sequence  collections.Sequence

	bankKeeper    bankkeeper.Keeper
	mailboxKeeper *mailboxkeeper.Keeper
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService storetypes.KVStoreService,
	authority string,
	bankKeeper bankkeeper.Keeper,
	mailboxKeeper *mailboxkeeper.Keeper,
) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:           cdc,
		addressCodec:  addressCodec,
		authority:     authority,
		HypTokens:     collections.NewMap(sb, types.HypTokenKey, "hyptokens", collections.BytesKey, codec.CollValue[types.HypToken](cdc)),
		Params:        collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Sequence:      collections.NewSequence(sb, types.HypTokensCountKey, "hyptokens_count"),
		bankKeeper:    bankKeeper,
		mailboxKeeper: mailboxKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

func (k Keeper) Handle(ctx context.Context, origin uint32, sender util.HexAddress, body []byte) error {

	// TODO implement
	fmt.Println("Message Received")
	fmt.Println(sender.String())
	fmt.Println(hex.EncodeToString(body))

	return nil
}
