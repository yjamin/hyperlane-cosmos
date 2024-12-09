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

func ParseWarpPayload(payload []byte) (WarpPayload, error) {
	if len(payload) != 64 {
		return WarpPayload{}, errors.New("payload is invalid")
	}

	amount := big.NewInt(0).SetBytes(payload[32:])

	return WarpPayload{
		recipient: payload[0:32],
		amount:    *amount,
	}, nil
}

func (p WarpPayload) Recipient() []byte {
	return p.recipient
}

func (p WarpPayload) Amount() *big.Int {
	newInt := big.NewInt(0)
	newInt.Set(&p.amount)
	return newInt
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
