package controller

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/lijie-keith/go_init_project/commonUtils"
	"github.com/lijie-keith/go_init_project/config"
	"github.com/lijie-keith/go_init_project/controller/model"
	"github.com/lijie-keith/go_init_project/store"
	"golang.org/x/crypto/sha3"
	"math/big"
	"net/http"
	"strconv"
)

func QueryAccountBalance(c *gin.Context) {
	addressStr := c.Query("address")

	address := common.HexToAddress(addressStr)

	balance, err := commonUtils.Client.BalanceAt(context.Background(), address, nil)

	if err != nil {
		config.SystemLogger.Error(err.Error())
	}
	c.JSON(http.StatusOK, commonUtils.OK.WithData(balance))
}

// QueryTokenBalance todo 未调通
func QueryTokenBalance(c *gin.Context) {
	addressStr := c.Query("address")
	tokenAddressStr := c.Query("tokenAddress")
	tokenAddress := common.HexToAddress(tokenAddressStr)
	instance, err := store.NewToken(tokenAddress, commonUtils.Client)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	address := common.HexToAddress(addressStr)
	of, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}
	c.JSON(http.StatusOK, commonUtils.OK.WithData(map[string]string{
		"name":    name,
		"balance": of.String(),
	}))
}

// CreateNewWallet 创建新的钱包
func CreateNewWallet(c *gin.Context) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}
	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		config.SystemLogger.Error("error casting public key")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	address = hexutil.Encode(hash.Sum(nil)[12:])
	c.JSON(http.StatusOK, commonUtils.OK.WithData(address))
}

func QueryLastBlock(c *gin.Context) {
	header, err := commonUtils.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	blockNumber := big.NewInt(header.Number.Int64())
	block, err := commonUtils.Client.BlockByNumber(context.Background(), blockNumber)
	fmt.Println(block.Number().Uint64())     // 5671744
	fmt.Println(block.Difficulty().Uint64()) // 3217000136609065
	fmt.Println(block.Hash().Hex())          // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9
	fmt.Println(len(block.Transactions()))   // 144
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}
	c.JSON(http.StatusOK, commonUtils.OK.WithData(block))
}

func QueryFirstTransactionFromBlock(c *gin.Context) {
	blockHeaderNumberStr := c.Query("headerNumber")
	atoi, err := strconv.Atoi(blockHeaderNumberStr)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}
	blockHeaderNumber := big.NewInt(int64(atoi))
	block, err := commonUtils.Client.BlockByNumber(context.Background(), blockHeaderNumber)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	firstTransactionHash := block.Transactions()[0].Hash()
	receipt, err := commonUtils.Client.TransactionReceipt(context.Background(), firstTransactionHash)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	c.JSON(http.StatusOK, commonUtils.OK.WithData(map[string]interface{}{
		"status": receipt.Status,
		"Logs":   receipt.Logs,
		"txHash": firstTransactionHash.Hex(),
	}))
}

func TransferBalance(c *gin.Context) {
	var transferRequest model.TransferRequest
	err := c.ShouldBind(&transferRequest)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	// 2.根据秘钥生成公钥地址和钱包地址
	privateKey, fromAddress, err := commonUtils.Encrypt(transferRequest.PrivateKey)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	nonce, err := commonUtils.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	value := big.NewInt(transferRequest.Amount)
	gasLimit := uint64(21000)
	gasPrice, err := commonUtils.Client.SuggestGasPrice(context.Background())

	toAddress := common.HexToAddress(transferRequest.ToAddress)

	chainID, err := commonUtils.Client.NetworkID(context.Background())
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	var data []byte

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     data,
	})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	err = commonUtils.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	c.JSON(http.StatusOK, commonUtils.OK.WithData(signedTx.Hash()))
}

func TransferTokenBalance(c *gin.Context) {
	// 1.绑定前端传参
	var transferRequest model.TransferRequest
	err := c.ShouldBind(&transferRequest)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	// 2.根据秘钥生成公钥地址和钱包地址
	privateKey, fromAddress, err := commonUtils.Encrypt(transferRequest.PrivateKey)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	// 3.生成交易所需要的参数
	nonce, err := commonUtils.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	value := big.NewInt(0)
	gasPrice, err := commonUtils.Client.SuggestGasPrice(context.Background())
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	toAddress := common.HexToAddress(transferRequest.ToAddress)
	tokenAddress := common.HexToAddress(transferRequest.TokenAddress)

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	amount := big.NewInt(transferRequest.Amount)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit := uint64(50000)
	//本地模拟测试后写死
	//gasLimit, err := commonUtils.Client.EstimateGas(context.Background(), ethereum.CallMsg{
	//	To:   &toAddress,
	//	Data: data,
	//})
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	chainID, err := commonUtils.Client.NetworkID(context.Background())
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	// 4.创建交易
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &tokenAddress,
		Value:    value,
		Data:     data,
	})

	// 5.交易签名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	// 6.生成交易
	err = commonUtils.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}

	rawTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		config.SystemLogger.Error(err.Error())
	}
	rawTxHex := hex.EncodeToString(rawTxBytes)

	c.JSON(http.StatusOK, commonUtils.OK.WithData(map[string]interface{}{
		"rawTx": rawTxHex,
		"tx":    signedTx.Hash().Hex(),
	}))
}
