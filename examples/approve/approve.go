package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	zerox "github.com/Daniil675/fiolet-playground-backend"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
	"github.com/lmittmann/w3/w3types"
)

// Define global variables

var (
	// Private key for signing the transaction
	privateKeyHex = "YOUR_PRIVATE_KEY"

	// Wallet address for transaction
	walletAddress = "YOUR_WALLET_ADDRESS"

	//URL of node
	nodeURL = "https://rpc.ankr.com/bsc"

	nonce       uint64 // Nonce value of the wallet address
	chainID     uint64 // Chain ID of the network
	estimateGas uint64 //Gas of transaction with given data

	txHash common.Hash // Hash of the transaction
)

func main() {
	// Create a new client for interacting with the 0x API
	client := zerox.NewClient(
		"YOUR_AUTH_KEY",
		zerox.BinanceSmartChain)
	// Get the price quote for a token swap from the 0x Swap API
	quote, err := client.GetPrice(context.Background(), zerox.PriceRequest{
		BuyToken:  "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c", // Address of a smart contract for the token you want to buy
		SellToken: "0x55d398326f99059fF775485246999027B3197955", // Address of a smart contract for the token you want to sell
		BuyAmount: "100000000000000000",                         // Amount to buy
	})
	if err != nil {
		fmt.Printf("GetQuote error:%s", err)
	}

	// Connect to the RPC endpoint
	clientw3 := w3.MustDial(nodeURL)
	defer clientw3.Close()

	// Create a new function call for the "approve" function with the specified arguments
	var funcApprove = w3.MustNewFunc("approve(address,uint256)", "bool")

	// Calculate the maximum value for the allowance
	maxValue := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	maxValue.Sub(maxValue, big.NewInt(1))

	// Encode the function arguments
	data, err := GetEncodedArgs(funcApprove, quote.AllowanceTarget, maxValue)
	if err != nil {
		fmt.Printf("GetEncodedArgs error:%s", err)
		return
	}

	msg := w3types.Message{
		From:  w3.A(walletAddress),
		To:    w3.APtr(quote.SellTokenAddress.Hex()),
		Input: data,
	}

	// Get the nonce, chain ID and estimated gas for transaction from the blockchain
	err = clientw3.Call(
		eth.Nonce(w3.A(walletAddress), nil).Returns(&nonce),
		eth.ChainID().Returns(&chainID),
		eth.EstimateGas(&msg, nil).Returns(&estimateGas),
	)
	if err != nil {
		fmt.Printf("Call error:%s", err)
		return
	}

	// Convert the chain ID to a big.Int
	chainIDBigInt := big.NewInt(quote.ChainID)

	// Build the approval transaction
	tx, ok := BuildApprove(nonce, quote, data, estimateGas)
	if !ok {
		fmt.Printf("BuildTx error:%s", err)
		return
	}

	//Convert the private key from hexadecimal to ECDSA format
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		fmt.Printf("PrivateKeyEncrypting error:%s", err)
		return
	}

	// Sign the approval transaction
	approveTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainIDBigInt), privateKey)
	if err != nil {
		fmt.Printf("Sign swapTx error:%s", err)
		return
	}

	//Send the approval transaction to the blockchain
	err = clientw3.Call(
		eth.SendTx(approveTx).Returns(&txHash),
	)
	if err != nil {
		fmt.Printf("Send approveTx error:%s", err)
		return
	}

	fmt.Println("Transaction hash:", txHash)
}

// BuildApprove builds an approval transaction for the given nonce and price quote.
// It returns the built transaction or an error if there was an issue.
func BuildApprove(nonce uint64, quote zerox.PriceResponse, inputData []byte, estimatedGas uint64) (*types.Transaction, bool) {
	// Convert the gas price from the price quote to a big.Int
	gasPrice, ok := GetBigInt(quote.GasPrice)
	if !ok {
		return nil, false
	}

	// Create the approval transaction
	approveTx := types.NewTransaction(
		nonce,
		quote.SellTokenAddress,
		big.NewInt(0),
		estimatedGas,
		gasPrice,
		inputData)

	return approveTx, true
}

// GetBigInt converts a string to a *big.Int. It returns the converted value and a boolean
// indicating whether the conversion was successful or not.
func GetBigInt(data string) (v *big.Int, ok bool) {
	v, ok = new(big.Int).SetString(data, 10)
	return
}

// GetByteFromString converts a hexadecimal string to a byte slice.
func GetByteFromString(data string) (v []byte, err error) {
	if len(data) > 2 {
		v, err = hex.DecodeString(data[2:])
	}
	return
}

// GetUint converts a string to an unsigned integer. It returns the converted value and any error that occurred.
func GetUint(data string) (gas uint64, err error) {
	gas, err = strconv.ParseUint(data, 10, 64)
	return
}

func GetEncodedArgs(function *w3.Func, args ...any) (input []byte, err error) {
	input, err = function.EncodeArgs(args...)
	return
}
