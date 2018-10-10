package blockchain

import (
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tyler-smith/go-bip39"
)

func derivePrivateKey(mnemonic string, path ...uint32) (*ecdsa.PrivateKey, error) {
	seed := bip39.NewSeed(mnemonic, "")
	curr, err := hdkeychain.NewMaster(seed, &chaincfg.Params{})
	if err != nil {
		return nil, err
	}
	for i, childIndex := range path {
		if i < 3 {
			childIndex += hdkeychain.HardenedKeyStart
		}
		curr, err = curr.Child(childIndex)
		if err != nil {
			return nil, err
		}
	}
	privKey, err := curr.ECPrivKey()
	if err != nil {
		return nil, err
	}
	return privKey.ToECDSA(), nil
}

func parseSignature(jobSignatureBytes []byte) (uint8, [32]byte, [32]byte, error) {
	r := [32]byte{}
	s := [32]byte{}

	if len(jobSignatureBytes) != 65 {
		return 0, r, s, fmt.Errorf("job signature incorrect length")
	}

	v := uint8(jobSignatureBytes[64])%27 + 27
	copy(r[:], jobSignatureBytes[0:32])
	copy(s[:], jobSignatureBytes[32:64])

	return v, r, s, nil
}

func addressToHex(address *common.Address) string {
	return address.Hex()
}

func bytesToBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func hexToBytes(str string) []byte {
	return common.FromHex(str)
}

func hexToAddress(str string) common.Address {
	return common.Address(common.BytesToAddress(hexToBytes(str)))
}
