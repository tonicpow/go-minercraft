package minercraft

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/libsv/go-bt/v2"
	"github.com/tonicpow/go-minercraft/v2/apis/mapi"
)

// FeeQuoteResponse is the raw response from the Merchant API request
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi#2-get-fee-quote
type FeeQuoteResponse struct {
	JSONEnvelope
	Quote *mapi.FeePayload `json:"quote"` // Custom field for unmarshalled payload data
}

// FeeQuote will fire a Merchant API request to retrieve the fees from a given miner
//
// This endpoint is used to get the different fees quoted by a miner.
// It returns a JSONEnvelope with a payload that contains the fees charged by a specific BSV miner.
// The purpose of the envelope is to ensure strict consistency in the message content for the purpose of signing responses.
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi#2-get-fee-quote
func (c *Client) FeeQuote(ctx context.Context, miner *Miner) (*FeeQuoteResponse, error) {

	// Make sure we have a valid miner
	if miner == nil {
		return nil, errors.New("miner was nil")
	}

	// Make the HTTP request
	result := getQuote(ctx, c, miner, mAPIRouteFeeQuote)
	if result.Response.Error != nil {
		return nil, result.Response.Error
	}

	// Parse the response
	response, err := result.parseFeeQuote()
	if err != nil {
		return nil, err
	}

	// Valid?
	if response.Quote == nil || len(response.Quote.Fees) == 0 {
		return nil, errors.New("failed getting quotes from: " + miner.Name)
	}

	isValid, err := response.IsValid()
	if err != nil {
		return nil, err
	}

	response.Validated = isValid

	// Return the fully parsed response
	return &response, nil
}

// internalResult is a shim for storing miner & http response data
type internalResult struct {
	Response *RequestResponse
	Miner    *Miner
}

// parseFeeQuote will convert the HTTP response into a struct and also unmarshal the payload JSON data
func (i *internalResult) parseFeeQuote() (response FeeQuoteResponse, err error) {

	// Process the initial response payload
	if err = response.process(i.Miner, i.Response.BodyContents); err != nil {
		return
	}

	// If we have a valid payload
	if len(response.Payload) > 0 {

		// Create a raw payload shim
		p := new(mapi.RawFeePayload)
		if err = json.Unmarshal([]byte(response.Payload), &p); err != nil {
			return
		}
		if response.Quote == nil {
			response.Quote = new(mapi.FeePayload)
		}

		// Create the response payload
		rawPayloadIntoQuote(p, response.Quote)
	}
	return
}

// rawPayloadIntoQuote will convert the raw parsed payload into a final quote payload
func rawPayloadIntoQuote(payload *mapi.RawFeePayload, quote *mapi.FeePayload) {

	// Set the fields from the raw payload into the quote
	quote.MinerID = payload.MinerID
	quote.APIVersion = payload.APIVersion
	quote.Timestamp = payload.Timestamp
	quote.ExpirationTime = payload.ExpirationTime
	quote.CurrentHighestBlockHash = payload.CurrentHighestBlockHash
	quote.CurrentHighestBlockHeight = payload.CurrentHighestBlockHeight
	quote.MinerReputation = payload.MinerReputation

	// Convert the mAPI fees into go-bt fees
	for _, f := range payload.Fees {
		t := bt.FeeTypeStandard
		if f.FeeType == mapi.FeeTypeData {
			t = bt.FeeTypeData
		}
		quote.Fees = append(quote.Fees, &bt.Fee{
			FeeType: t,
			MiningFee: bt.FeeUnit{
				Satoshis: f.MiningFee.Satoshis,
				Bytes:    f.MiningFee.Bytes,
			},
			RelayFee: bt.FeeUnit{
				Satoshis: f.RelayFee.Satoshis,
				Bytes:    f.RelayFee.Bytes,
			},
		})
	}
}

// getQuote will fire the HTTP request to retrieve the fee/policy quote
func getQuote(ctx context.Context, client *Client, miner *Miner, route string) (result *internalResult) {
	sb := strings.Builder{}

	api, err := client.MinerAPIByMinerID(miner.MinerID, client.apiType)
	if err != nil {
		result.Response = &RequestResponse{Error: err}
		return
	}

	sb.WriteString(api.URL + route)
	result = &internalResult{Miner: miner}
	quoteURL, err := url.Parse(sb.String())
	if err != nil {
		result.Response = &RequestResponse{Error: err}
		return
	}

	result = &internalResult{Miner: miner}
	result.Response = httpRequest(ctx, client, &httpPayload{
		Method: http.MethodGet,
		URL:    quoteURL.String(),
		Token:  api.Token,
	})
	return
}
