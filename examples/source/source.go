package main

import (
	"context"
	"fmt"

	zerox "github.com/Daniil675/fiolet-playground-backend"
)

func main() {
	// Initialize the 0x Swap API client
	client := zerox.NewClient(
		"YOUR_AUTH_KEY",
		zerox.BinanceSmartChain)
	// Getting a list of sources of liquidity for this network
	sources, err := client.GetSources(context.Background())
	if err != nil {
		fmt.Printf("GetSources error:%s", err)
	}

	fmt.Println(sources)
}
