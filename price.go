package zerox

import (
	"context"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-querystring/query"
)

type (
	PriceRequest struct {
		SellToken                       string `url:"sellToken"`
		BuyToken                        string `url:"buyToken"`
		SellAmount                      string `url:"sellAmount,omitempty"`
		BuyAmount                       string `url:"buyAmount,omitempty"`
		SlippagePercentage              string `url:"slippagePercentage,omitempty"`
		GasPrice                        string `url:"gasPrice,omitempty"`
		TakerAddress                    string `url:"takerAddress,omitempty"`
		ExcludedSources                 string `url:"excludedSources,omitempty"`
		IncludedSources                 string `url:"skipValidation,omitempty"`
		IntentOnFilling                 string `url:"intentOnFilling,omitempty"`
		FeeRecipient                    string `url:"feeRecipient,omitempty"`
		BuyTokenPercentageFee           string `url:"buyTokenPercentageFee,omitempty"`
		EnableSlippageProtection        string `url:"enableSlippageProtection,omitempty"`
		PriceImpactProtectionPercentage string `url:"priceImpactProtectionPercentage,omitempty"`
	}

	PriceResponse struct {
		ChainID              int64            `json:"chainId"`
		Price                string           `json:"price"`
		GrossPrice           string           `json:"grossPrice"`
		EstimatedPriceImpact string           `json:"estimatedPriceImpact"`
		Value                string           `json:"value"`
		GasPrice             string           `json:"gasPrice"`
		Gas                  string           `json:"gas"`
		EstimatedGas         string           `json:"estimatedGas"`
		ProtocolFee          string           `json:"protocolFee"`
		MinimumProtocolFee   string           `json:"minimumProtocolFee"`
		BuyTokenAddress      common.Address   `json:"buyTokenAddress"`
		BuyAmount            string           `json:"buyAmount"`
		GrossBuyAmount       string           `json:"grossBuyAmount"`
		SellTokenAddress     common.Address   `json:"sellTokenAddress"`
		SellAmount           string           `json:"sellAmount"`
		GrossSellAmount      string           `json:"grossSellAmount"`
		Sources              []Source         `json:"sources"`
		AllowanceTarget      common.Address   `json:"allowanceTarget"`
		SellTokenToEthRate   string           `json:"sellTokenToEthRate"`
		BuyTokenToEthRate    string           `json:"buyTokenToEthRate"`
		ExpectedSlippage     string           `json:"expectedSlippage"`
		AuxiliaryChainData   interface{}      `json:"auxiliaryChainData"`
		Fees                 map[string][]Fee `json:"fees"`
	}

	Fee struct {
		Name   string
		Amount string
	}
)

func (c *Client) GetPrice(
	ctx context.Context,
	request PriceRequest,
) (PriceResponse, error) {
	var response PriceResponse

	urlSuffix := "swap/v1/price"

	req, err := c.newRequest(ctx, http.MethodGet, c.config.URL+urlSuffix)
	if err != nil {
		return PriceResponse{}, err
	}

	parametrs, err := query.Values(request)
	if err != nil {
		return response, err
	}

	req.URL.RawQuery = parametrs.Encode()

	err = c.sendRequest(req, &response)
	if err != nil {
		return PriceResponse{}, err
	}
	return response, nil
}
