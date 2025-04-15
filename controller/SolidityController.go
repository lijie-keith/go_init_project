package controller

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/lijie-keith/go_init_project/commonUtils"
	"github.com/lijie-keith/go_init_project/config"
	"github.com/lijie-keith/go_init_project/controller/model"
	store "github.com/lijie-keith/go_init_project/store"
	"math/big"
	"net/http"
	"strings"
)

func DeploySolidity(c *gin.Context) {
	var transferRequest model.TransferRequest
	if err := c.ShouldBind(&transferRequest); err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	privateKey, fromAddress, err := commonUtils.Encrypt(transferRequest.PrivateKey)
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}

	nonce, err := commonUtils.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}

	gasPrice, err := commonUtils.Client.SuggestGasPrice(context.Background())
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}

	chainID, err := commonUtils.Client.NetworkID(context.Background())
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	auth := bind.NewKeyedTransactor(privateKey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	input := "1.0"
	address, tx, instance, err := store.DeployStore(auth, commonUtils.Client, input)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	_ = instance
	c.JSON(http.StatusOK, commonUtils.OK.WithData(map[string]interface{}{
		"address": address.Hex(),
		"tx":      tx.Hash().Hex(),
	}))

}

func GetSolidityByteCode(c *gin.Context) {
	var transferRequest model.TransferRequest
	if err := c.ShouldBind(&transferRequest); err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	address := common.HexToAddress(transferRequest.TokenAddress)
	bytecode, err := commonUtils.Client.CodeAt(context.Background(), address, nil)
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, commonUtils.OK.WithData(map[string]interface{}{
		"bytecode": hex.EncodeToString(bytecode),
	}))
}

func GetSolidityEventLogs(c *gin.Context) {
	var transferRequest model.TransferRequest
	if err := c.ShouldBind(&transferRequest); err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	address := common.HexToAddress(transferRequest.TokenAddress)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(8123902),
		ToBlock:   big.NewInt(8123902),
		Addresses: []common.Address{address},
	}

	logs, err := commonUtils.Client.FilterLogs(context.Background(), query)
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	contractAbi, err := abi.JSON(strings.NewReader(store.StoreMetaData.ABI))
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}

	for _, vLog := range logs {
		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}
		err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			config.SystemLogger.Error(err.Error())
			return
		}
		c.JSON(http.StatusOK, commonUtils.OK.WithData(map[string]interface{}{
			"key":   event.Key[:],
			"value": event.Value[:],
		}))
		var topics [4]string
		for i := range vLog.Topics {
			topics[i] = vLog.Topics[i].Hex()
		}

	}
	evenSignature := []byte("ItemSet[bytes32,bytes32]")
	hash := crypto.Keccak256(evenSignature)
	_ = hash

}

func NewStore(c *gin.Context) {
	var transferRequest model.TransferRequest
	if err := c.ShouldBind(&transferRequest); err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	address := common.HexToAddress(transferRequest.TokenAddress)
	instance, err := store.NewStore(address, commonUtils.Client)
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	version, err := instance.Version(nil)
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, commonUtils.OK.WithData(version))
}

func SetItem(c *gin.Context) {
	var transferRequest model.TransferRequest
	if err := c.ShouldBind(&transferRequest); err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	privateKey, fromAddress, err := commonUtils.Encrypt(transferRequest.PrivateKey)
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	nonce, err := commonUtils.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	gasPrice, err := commonUtils.Client.SuggestGasPrice(context.Background())
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}

	chainID, err := commonUtils.Client.NetworkID(context.Background())
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	auth := bind.NewKeyedTransactor(privateKey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	address := common.HexToAddress(transferRequest.TokenAddress)
	instance, err := store.NewStore(address, commonUtils.Client)
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], []byte("foo"))
	copy(value[:], []byte("bar"))
	tx, err := instance.SetItem(auth, key, value)
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	_ = tx
	items, err := instance.Items(nil, key)
	if err != nil {
		config.SystemLogger.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, commonUtils.OK.WithData(items[:]))
}
