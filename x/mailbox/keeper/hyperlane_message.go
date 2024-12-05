package keeper

import (
	"encoding/binary"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	VersionOffset     = 0
	NonceOffset       = 1
	OriginOffset      = 5
	SenderOffset      = 9
	DestinationOffset = 41
	RecipientOffset   = 45
	BodyOffset        = 77
)

func Id(message []byte) []byte {
	return crypto.Keccak256(message)
}

func Version(message []byte) byte {
	return message[VersionOffset]
}

func Nonce(message []byte) uint32 {
	return binary.BigEndian.Uint32(message[NonceOffset:OriginOffset])
}

func Origin(message []byte) uint32 {
	return binary.BigEndian.Uint32(message[OriginOffset:SenderOffset])
}

func Sender(message []byte) []byte {
	return message[SenderOffset:DestinationOffset]
}

func Destination(message []byte) uint32 {
	return binary.BigEndian.Uint32(message[DestinationOffset:RecipientOffset])
}

func Recipient(message []byte) []byte {
	return message[RecipientOffset:BodyOffset]
}

func Body(message []byte) []byte {
	return message[BodyOffset:]
}
