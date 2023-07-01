# zerox

<img src="https://github.com/fioletxyz/0x-go/assets/81973532/2c28f92e-d4b3-4781-88fd-5c618caf149d" align="right" alt="W3 Gopher" width="158" height="224">

Package `zerox` is a wrapped 0x Swap API.

This package has a respresentation of:

- **Quote** **(swap/v1/quote)** - Get an easy-to-consume quote for buying or selling any ERC20 token. Returns a transaction that can be submitted to an Ethereum node.
- **Price** **(swap/v1/price)** - /price is nearly identical to /quote, but with a few key differences. /price does not return a transaction that can be submitted on-chain; it simply provides us the same information. Think of it as the "read-only" version of /quote.
- **Source** **(swap/v1/source)** - Returns the liquidity sources enabled for the chain.

## Install

```
go get github.com/fioletxyz/0x-go
```

## Getting Started

> **Note**
> Check out the [examples](examples/)!

Connect your API Auth key and choose an endpoint.

```go
// Create a new client for interacting with the 0x API
client := zerox.NewClient(
	"YOUR_AUTH_KEY",
	zerox.BinanceSmartChain) //Choose an endpoint from config.go
```
