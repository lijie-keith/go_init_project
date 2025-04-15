package commonUtils

import (
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lijie-keith/go_init_project/config"
)

func Encrypt(privateKey string) (*ecdsa.PrivateKey, common.Address, error) {
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return nil, common.Address{}, err
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		config.SystemLogger.Error("error casting public key")
		return nil, common.Address{}, errors.New("error casting public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKeyECDSA, fromAddress, nil
}
