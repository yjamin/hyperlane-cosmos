package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide the data as an argument")
	}
	dataHex := os.Args[1]

	if strings.HasPrefix(dataHex, "0x") {
		dataHex = dataHex[2:]
	}
	data, err := hex.DecodeString(dataHex)
	if err != nil {
		panic(err)
	}

	hash := crypto.Keccak256Hash(data)

	// Include all private keys of the used ISM.
	privateKeys := []string{
		// PubKey: 0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05
		"fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a18",
		// PubKey: 0x0417f57017d748288ccf6341993e47618ce3d4d60614ae09f5149acec191fad3fbca5a8ce4144077948c843ea8e863e3997b6da7a1a6d6c9708f658371430ce06b
		"c87509a1c067bbde78beb793e6fa49e3462d7e7bcfbb4e3f79e926fc27ae42c4",
		// PubKey: 0x04ce7edc292d7b747fab2f23584bbafaffde5c8ff17cf689969614441e0527b90015ea9fee96aed6d9c0fc2fbe0bd1883dee223b3200246ff1e21976bdbc9a0fc8
		"ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f",
	}

	var signatures []byte

	for i, pk := range privateKeys {
		privateKey, err := crypto.HexToECDSA(pk)
		if err != nil {
			log.Fatal(err)
		}

		publicKey := privateKey.Public()
		_, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}

		signature, err := crypto.Sign(hash.Bytes(), privateKey)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Signature %d: %s\n", i+1, hexutil.Encode(signature))
		signatures = append(signatures, signature...)
	}

	signaturesHex := hexutil.Encode(signatures)
	fmt.Println("Metadata:", signaturesHex)

	fmt.Println(fmt.Sprintf("\n./build/kyved tx mailbox process %s %s", signaturesHex, dataHex))
}
