package crypto

import (
	"crypto/ecdsa"
	"fmt"
)

func FromMnemonic(mnemonic string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	wallet, err := NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, nil, err
	}

	path := MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, nil, err
	}

	fmt.Printf("Account address: %s\n", account.Address.Hex())

	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := wallet.PublicKey(account)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}
