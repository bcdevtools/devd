package utils

import (
	"crypto/ecdsa"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	cosmoshd "github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/accounts"
)

func FromMnemonicToPrivateKey(mnemonic, password string) (*ecdsa.PrivateKey, error) {
	hdPathStr := cosmoshd.CreateHDPath(60, 0, 0).String()
	hdPath, err := accounts.ParseDerivationPath(hdPathStr)
	if err != nil {
		return nil, err
	}

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		return nil, err
	}

	// create a BTC-utils hd-derivation keychain
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	key := masterKey
	for _, n := range hdPath {
		key, err = key.Derive(n)
		if err != nil {
			return nil, err
		}
	}

	// btc-utils representation of a secp256k1 private key
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}

	// cast private key to a convertible form (single scalar field element of secp256k1)
	// and then load into ethcrypto private key format.
	return privateKey.ToECDSA(), nil
}
