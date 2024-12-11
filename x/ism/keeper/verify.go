package keeper

import (
	"bytes"
	"context"
	"github.com/KYVENetwork/hyperlane-cosmos/util"
	"github.com/KYVENetwork/hyperlane-cosmos/x/ism/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func (k Keeper) Verify(ctx context.Context, ismId string, rawMetadata []byte, messageStr string) (bool, error) {
	// Retrieve ISM
	ism, err := k.Isms.Get(ctx, ismId)
	if err != nil {
		return false, err
	}

	message, err := util.DecodeEthHex(messageStr)
	if err != nil {
		return false, err
	}

	hash := crypto.Keccak256Hash(message)

	// Get ISM validator public keys
	var validatorPubKeys [][]byte
	for _, pubKeyStr := range ism.ValidatorPubKeys {
		pubKey, err := util.DecodeEthHex(pubKeyStr)
		if err != nil {
			return false, err
		}
		validatorPubKeys = append(validatorPubKeys, pubKey)
	}

	// Get signature count
	numSignatures := (len(rawMetadata) - types.SIGNATURES_OFFSET) / types.SIGNATURE_LENGTH

	validCount := 0
	for i := uint32(0); i < uint32(numSignatures); i++ {
		sig := types.SignatureAt(rawMetadata, i)

		recoveredPubKey, err := crypto.SigToPub(hash.Bytes(), sig)
		if err != nil {
			return false, err
		}

		for _, validatorPubKey := range validatorPubKeys {
			if bytes.Equal(crypto.FromECDSAPub(recoveredPubKey), validatorPubKey) {
				validCount++
				break
			}
		}
	}

	if validCount >= int(ism.Threshold) {
		return true, nil
	}
	return false, nil
}

func (k Keeper) IsmIdExists(ctx context.Context, ismId string) (bool, error) {
	ism, err := k.Isms.Has(ctx, ismId)
	if err != nil {
		return false, err
	}
	return ism, nil
}
