package zerox

import (
	"context"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-querystring/query"
)

type (
	QuoteRequest struct {
		SellToken                       string `url:"sellToken"`
		BuyToken                        string `url:"buyToken"`
		SellAmount                      string `url:"sellAmount,omitempty"`
		BuyAmount                       string `url:"buyAmount,omitempty"`
		SlippagePercentage              string `url:"slippagePercentage,omitempty"`
		GasPrice                        string `url:"gasPrice,omitempty"`
		TakerAddress                    string `url:"takerAddress,omitempty"`
		ExcludedSources                 string `url:"excludedSources,omitempty"`
		IncludedSources                 string `url:"includedSources,omitempty"`
		SkipValidation                  string `url:"skipValidation,omitempty"`
		IntentOnFilling                 string `url:"intentOnFilling,omitempty"`
		FeeRecipient                    string `url:"feeRecipient,omitempty"`
		BuyTokenPercentageFee           string `url:"buyTokenPercentageFee,omitempty"`
		EnableSlippageProtection        string `url:"enableSlippageProtection,omitempty"`
		PriceImpactProtectionPercentage string `url:"priceImpactProtectionPercentage,omitempty"`
	}

	QuoteResponse struct {
		ChainID              int64          `json:"chainId"`
		Price                string         `json:"price"`
		GrossPrice           string         `json:"grossPrice"`
		GuaranteedPrice      string         `json:"guaranteedPrice"`
		EstimatedPriceImpact string         `json:"estimatedPriceImpact"`
		To                   string         `json:"to"`
		Data                 string         `json:"data"`
		Value                string         `json:"value"`
		GasPrice             string         `json:"gasPrice"`
		Gas                  string         `json:"gas"`
		EstimatedGas         string         `json:"estimatedGas"`
		ProtocolFee          string         `json:"protocolFee"`
		MinimumProtocolFee   string         `json:"minimumProtocolFee"`
		BuyAmount            string         `json:"buyAmount"`
		GrossBuyAmount       string         `json:"grossBuyAmount"`
		SellAmount           string         `json:"sellAmount"`
		GrossSellAmount      string         `json:"grossSellAmount"`
		Sources              []Source       `json:"sources"`
		BuyTokenAddress      common.Address `json:"buyTokenAddress"`
		SellTokenAddress     common.Address `json:"sellTokenAddress"`
		AllowanceTarget      common.Address `json:"allowanceTarget"`
		Orders               []Order        `json:"orders"`
		SellTokenToEthRate   string         `json:"sellTokenToEthRate"`
		BuyTokenToEthRate    string         `json:"buyTokenToEthRate"`
		ExpectedSlippage     string         `json:"expectedSlippage"`
	}

	Source struct {
		Name       string `json:"name"`
		Proportion string `json:"proportion"`
	}

	Order struct {
		Type          int      `json:"type"`
		SourceOfOrder string   `json:"source"`
		MakerToken    string   `json:"makerToken"`
		TakerToken    string   `json:"takerToken"`
		MakerAmount   string   `json:"makerAmount"`
		TakerAmount   string   `json:"takerAmount"`
		FillData      FillData `json:"fillData"`
		Fill          Fill     `json:"fill"`
	}

	FillData struct {
		TokenAddressPath []string `json:"tokenAddressPath,omitempty"`
		Path             string   `json:"path,omitempty"`
		Router           string   `json:"router"`
		GasUsed          string   `json:"gasUsed,omitempty"`
	}

	Fill struct {
		Input          string `json:"input"`
		Output         string `json:"output"`
		AdjustedOutput string `json:"adjustedOutput"`
		Gas            int    `json:"gas"`
	}
)

func (c *Client) GetQuote(
	ctx context.Context,
	request QuoteRequest,
) (QuoteResponse, error) {
	var response QuoteResponse

	urlSuffix := "swap/v1/quote"

	req, err := c.newRequest(ctx, http.MethodGet, c.config.URL+urlSuffix)
	if err != nil {
		return QuoteResponse{}, err
	}

	parametrs, err := query.Values(request)
	if err != nil {
		return response, err
	}

	req.URL.RawQuery = parametrs.Encode()

	err = c.sendRequest(req, &response)
	if err != nil {
		return QuoteResponse{}, err
	}
	return response, nil
}
