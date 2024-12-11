package keeper

import (
	"context"
	"github.com/KYVENetwork/hyperlane-cosmos/util"
	"github.com/KYVENetwork/hyperlane-cosmos/x/ism/types"
)

func (ms msgServer) CreateMultisigIsm(ctx context.Context, req *types.MsgCreateMultisigIsm) (*types.MsgCreateMultisigIsmResponse, error) {
	ismCount, err := ms.k.IsmsSequence.Next(ctx)
	if err != nil {
		return nil, err
	}

	prefixedId := util.CreateHexAddress(types.ModuleName, int64(ismCount))

	newIsm := types.MultiSigIsm{
		Creator:          req.Creator,
		ValidatorPubKeys: req.MultiSig.ValidatorPubKeys,
		Threshold:        req.MultiSig.Threshold,
		Id:               prefixedId.String(),
		IsmType:          5,
	}

	if err = ms.k.Isms.Set(ctx, prefixedId.String(), newIsm); err != nil {
		return nil, err
	}

	return &types.MsgCreateMultisigIsmResponse{}, nil
}
