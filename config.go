package zerox

import "net/http"

type (
	Config struct {
		AuthToken string
		URL       string

		HTTPClient *http.Client
	}
)

const (
	EthereumMainnet   = "https://api.0x.org/"
	EthereumGoerli    = "https://goerli.api.0x.org/"
	BinanceSmartChain = "https://bsc.api.0x.org/"
	Polygon           = "https://polygon.api.0x.org/"
	PolygonMumbai     = "https://mumbai.api.0x.org"
	Optimism          = "https://optimism.api.0x.org/"
	Fantom            = "https://fantom.api.0x.org/"
	Celo              = "https://celo.api.0x.org/"
	Avalanche         = "https://avalanche.api.0x.org/"
	Arbitrum          = "https://arbitrum.api.0x.org/"
)

func GetConfig(authToken string, url string) Config {
	return Config{
		AuthToken:  authToken,
		URL:        url,
		HTTPClient: &http.Client{},
	}
}
