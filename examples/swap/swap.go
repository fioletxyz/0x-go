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
)

var (
	// Private key for signing the transaction
	privateKeyHex = "YOUR_PRIVATE_KEY"

	// Wallet address for transaction
	walletAddress = "YOUR_WALLET_ADDRESS"

	nonce   uint64 // Nonce value of the wallet address
	chainID uint64 // Chain ID of the network

	//URL of node
	nodeURL = "https://rpc.ankr.com/bsc"

	txHash common.Hash
)

func main() {
	// Initialize the 0x Swap API client
	client := zerox.NewClient(
		"YOUR_AUTH_KEY",         // 0x Swap API key
		zerox.BinanceSmartChain, // Network
	)

	// Get a quote for the trade
	quote, err := client.GetQuote(context.Background(), zerox.QuoteRequest{
		BuyToken:  "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c", // Address of a smart contract for the token you want to buy
		SellToken: "0x55d398326f99059fF775485246999027B3197955", // Address of a smart contract for the token you want to sell
		BuyAmount: "100000000000000000",                         // Amount to buy
	})
	if err != nil {
		fmt.Printf("GetQuote error: %s", err)
	}

	// Connect to the node
	clientw3 := w3.MustDial(nodeURL)
	defer clientw3.Close()

	// Get the wallet's nonce and chain ID from the Ethereum node
	err = clientw3.Call(
		eth.Nonce(w3.A(walletAddress), nil).Returns(&nonce), // Get the nonce (transaction count) for the wallet address
		eth.ChainID().Returns(&chainID),                     // Get the chain ID of the network
	)
	if err != nil {
		fmt.Printf("Call error: %s", err)
	}

	chainID := big.NewInt(quote.ChainID) // Create a new big.Int from the chain ID

	// Build the transaction
	tx, err := BuildTx(nonce, quote)
	if err != nil {
		fmt.Printf("BuildTx error: %s", err)
	}

	// Convert the private key from hex to ECDSA format
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		fmt.Printf("PrivateKeyEncrypting error: %s", err)
	}

	// Sign the transaction using the private key
	tx, err = types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		fmt.Printf("Sign swapTx error: %s", err)
	}

	// Send the signed transaction to the node
	err = clientw3.Call(
		eth.SendTx(tx).Returns(&txHash), // Send the transaction to the network
	)
	if err != nil {
		fmt.Printf("Send swapTx error: %s", err)
	}

	fmt.Println("Transaction hash:", txHash)
}

// BuildTx builds a new transaction using the provided nonce and quote information
func BuildTx(nonce uint64, quote zerox.QuoteResponse) (*types.Transaction, error) {
	maxValue := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	maxValue.Sub(maxValue, big.NewInt(1))

	estimatedGas, err := GetUint(quote.EstimatedGas)
	if err != nil {
		fmt.Printf("BuildTx error: %s", err)
	}

	gasPrice, ok := GetBigInt(quote.GasPrice)
	if !ok {
		fmt.Printf("BuildTx error: %s", err)
	}

	data, err := GetByteFromString(quote.Data)
	if err != nil {
		fmt.Printf("BuildTx error: %s", err)
	}

	// Create a new transaction object
	tx := types.NewTransaction(
		nonce,                 // Nonce of the transaction
		quote.AllowanceTarget, // Target address of the transaction
		big.NewInt(0),         // Value (amount) of the transaction
		estimatedGas,          // Estimated gas for the transaction
		gasPrice,              // Gas price for the transaction
		data,                  // Data payload of the transaction
	)

	return tx, nil
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
