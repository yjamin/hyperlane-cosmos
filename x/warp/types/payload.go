package types

import (
	"errors"
	"math/big"
	"slices"
)

type WarpPayload struct {
	recipient []byte
	amount    big.Int
}

func NewWarpPayload(recipient []byte, amount big.Int) (WarpPayload, error) {
	if len(amount.Bytes()) > 32 {
		return WarpPayload{}, errors.New("amount is too long")
	}
	if len(recipient) > 32 {
		return WarpPayload{}, errors.New("recipient address is too long")
	}

	return WarpPayload{recipient: recipient, amount: amount}, nil
}

func (p WarpPayload) Bytes() []byte {

	intBytes := p.amount.Bytes()
	amountBytes := make([]byte, 32)
	copy(amountBytes[32-len(intBytes):], intBytes)

	recBytes := p.recipient
	receiverBytes := make([]byte, 32)
	copy(receiverBytes[32-len(recBytes):], recBytes)

	return slices.Concat(
		receiverBytes,
		amountBytes,
	)
}
