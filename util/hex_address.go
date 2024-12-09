package util

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type HexAddress [32]byte

func (h HexAddress) String() string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(h[:]))
}

func (h HexAddress) Bytes() []byte {
	return h[:]
}

func DecodeHexAddress(s string) (HexAddress, error) {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}

	if len(s) != 64 {
		return HexAddress{}, errors.New("invalid hex address length")
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return HexAddress{}, err
	}
	return HexAddress(b), nil
}

func DecodeEthHex(s string) ([]byte, error) {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func CreateHexAddress(identifier string, id int64) HexAddress {
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, uint64(id))
	message := append([]byte(identifier), idBytes...)
	return sha256.Sum256(message)
}
