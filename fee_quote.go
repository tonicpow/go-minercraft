package minercraft

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const (

	// FeeTypeData is the key corresponding to the data rate
	FeeTypeData = "data"

	// FeeTypeStandard is the key corresponding to the standard rate
	FeeTypeStandard = "standard"

	// FeeCategoryMining is the category corresponding to the mining rate
	FeeCategoryMining = "mining"

	// FeeCategoryRelay is the category corresponding to the relay rate
	FeeCategoryRelay = "relay"
)

/*
Example feeQuote response from Merchant API:

{
	"payload": "{\"apiVersion\":\"0.1.0\",\"timestamp\":\"2020-10-07T21:13:04.335Z\",\"expiryTime\":\"2020-10-07T21:23:04.335Z\",\"minerId\":\"0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087\",\"currentHighestBlockHash\":\"000000000000000000edb30c3bbbc8e6a07e522e85522e6a213f7e933e6e2d8d\",\"currentHighestBlockHeight\":655874,\"minerReputation\":null,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}}]}",
	"signature": "304402206443bea5bdd98a16e23eb61c36b4b998bd68ceb9c84983c7e695e267b21a30440220191571e9b9632c8337d9196723ca20eefa63966ef6360170db0e57a04047453f",
	"publicKey": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087",
	"encoding": "UTF-8",
	"mimetype": "application/json"
}
*/

// FeeQuoteResponse is the raw response from the Merchant API request
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi/tree/v1.2-beta#get-fee-quote
type FeeQuoteResponse struct {
	JSONEnvelope
	Quote *FeePayload `json:"quote"` // Custom field for unmarshalled payload data
}

/*
Example FeeQuoteResponse.Payload (unmarshalled):

{
  "apiVersion": "0.1.0",
  "timestamp": "2020-10-07T21:13:04.335Z",
  "expiryTime": "2020-10-07T21:23:04.335Z",
  "minerId": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087",
  "currentHighestBlockHash": "000000000000000000edb30c3bbbc8e6a07e522e85522e6a213f7e933e6e2d8d",
  "currentHighestBlockHeight": 655874,
  "minerReputation": null,
  "fees": [
    {
      "feeType": "standard",
      "miningFee": {
        "satoshis": 500,
        "bytes": 1000
      },
      "relayFee": {
        "satoshis": 250,
        "bytes": 1000
      }
    },
    {
      "feeType": "data",
      "miningFee": {
        "satoshis": 500,
        "bytes": 1000
      },
      "relayFee": {
        "satoshis": 250,
        "bytes": 1000
      }
    }
  ]
}
*/

// FeePayload is the unmarshalled version of the payload envelope
type FeePayload struct {
	APIVersion                string      `json:"apiVersion"`
	Timestamp                 string      `json:"timestamp"`
	ExpirationTime            string      `json:"expiryTime"`
	MinerID                   string      `json:"minerId"`
	CurrentHighestBlockHash   string      `json:"currentHighestBlockHash"`
	CurrentHighestBlockHeight uint64      `json:"currentHighestBlockHeight"`
	MinerReputation           interface{} `json:"minerReputation"` // Not sure what this value is
	Fees                      []*feeType  `json:"fees"`
}

// CalculateFee will return the fee for the given txBytes
// Type: "FeeTypeData" or "FeeTypeStandard"
// Category: "FeeCategoryMining" or "FeeCategoryRelay"
//
// If no fee is found or fee is 0, returns 1 & error
//
// Spec: https://github.com/bitcoin-sv-specs/brfc-misc/tree/master/feespec#deterministic-transaction-fee-calculation-dtfc
func (f *FeePayload) CalculateFee(feeCategory, feeType string, txBytes uint64) (uint64, error) {

	// Valid feeType?
	if !strings.EqualFold(feeType, FeeTypeData) && !strings.EqualFold(feeType, FeeTypeStandard) {
		return 0, fmt.Errorf("feeType %s is not recognized", feeType)
	} else if !strings.EqualFold(feeCategory, FeeCategoryMining) && !strings.EqualFold(feeCategory, FeeCategoryRelay) {
		return 0, fmt.Errorf("feeCategory %s is not recognized", feeCategory)
	}

	// Loop all fee types looking for feeType (data or standard)
	for _, fee := range f.Fees {

		// Detect the type (data or standard)
		if fee.FeeType == feeType {

			// Multiply & Divide
			var calcFee uint64
			if strings.EqualFold(feeCategory, FeeCategoryMining) {
				calcFee = (fee.MiningFee.Satoshis * txBytes) / fee.MiningFee.Bytes
			} else {
				calcFee = (fee.RelayFee.Satoshis * txBytes) / fee.RelayFee.Bytes
			}

			// Check for zero
			if calcFee != 0 {
				return calcFee, nil
			}

			// If txBytes is zero this error will occur
			return 1, fmt.Errorf("warning: fee calculation was 0")
		}
	}

	// No fee type found in the slice of fees
	return 1, fmt.Errorf("feeType %s is not found in fees", feeType)
}

/*
Example FeePayload.Fees type:
{
  "feeType": "standard",
  "miningFee": {
	"satoshis": 500,
	"bytes": 1000
  },
  "relayFee": {
	"satoshis": 250,
	"bytes": 1000
  }
}
*/

// feeType is the the corresponding type of fee (standard or data)
type feeType struct {
	FeeType   string     `json:"feeType"`
	MiningFee *feeAmount `json:"miningFee"`
	RelayFee  *feeAmount `json:"relayFee"`
}

// feeAmount is the actual fee for the given feeType
type feeAmount struct {
	Bytes    uint64 `json:"bytes"`
	Satoshis uint64 `json:"satoshis"`
}

// FeeQuote will fire a Merchant API request to retrieve the fees from a given miner
//
// This endpoint is used to get the different fees quoted by a miner.
// It returns a JSONEnvelope with a payload that contains the fees charged by a specific BSV miner.
// The purpose of the envelope is to ensure strict consistency in the message content for the purpose of signing responses.
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi/tree/v1.2-beta#get-fee-quote
func (c *Client) FeeQuote(miner *Miner) (*FeeQuoteResponse, error) {

	// Make sure we have a valid miner
	if miner == nil {
		return nil, errors.New("miner was nil")
	}

	// Make the HTTP request
	result := getQuote(c, miner)
	if result.Response.Error != nil {
		return nil, result.Response.Error
	}

	// Parse the response
	response, err := result.parseQuote()
	if err != nil {
		return nil, err
	}

	// Valid?
	if response.Quote == nil || len(response.Quote.Fees) == 0 {
		return nil, errors.New("failed getting quotes from: " + miner.Name)
	}

	// Return the fully parsed response
	return &response, nil
}

// BestQuote will check all known miners and compare rates, returning the best rate/quote
//
// Note: this might return different results each time if miners have the same rates as
// it's a race condition on which results come back first
func (c *Client) BestQuote(feeCategory, feeType string) (*FeeQuoteResponse, error) {

	// Best rate & quote
	var bestRate uint64
	var bestQuote FeeQuoteResponse

	// The channel for the internal results
	resultsChannel := make(chan *internalResult, len(c.Miners))

	// Loop each miner (break into a Go routine for each quote request)
	var wg sync.WaitGroup
	for _, miner := range c.Miners {
		wg.Add(1)
		go getQuoteRoutine(&wg, c, miner, resultsChannel)
	}

	// Waiting for all requests to finish
	wg.Wait()
	close(resultsChannel)

	// Loop the results of the channel
	var testRate uint64
	for result := range resultsChannel {

		// Check for error?
		if result.Response.Error != nil {
			return nil, result.Response.Error
		}

		// Parse the response
		quote, err := result.parseQuote()
		if err != nil {
			return nil, err
		}

		// Do we have a rate set?
		if bestRate == 0 {
			bestQuote = quote
			if bestRate, err = quote.Quote.CalculateFee(feeCategory, feeType, 1000); err != nil {
				return nil, err
			}
		} else { // Test the other quotes
			if testRate, err = quote.Quote.CalculateFee(feeCategory, feeType, 1000); err != nil {
				return nil, err
			}
			if testRate < bestRate {
				bestRate = testRate
				bestQuote = quote
			}
		}
	}

	// Return the best quote found
	return &bestQuote, nil
}

// internalResult is a shim for storing miner & http response data
type internalResult struct {
	Response *RequestResponse
	Miner    *Miner
}

// parseQuote will convert the HTTP response into a struct and also unmarshal the payload JSON data
func (i *internalResult) parseQuote() (response FeeQuoteResponse, err error) {

	// Process the initial response payload
	if err = response.process(i.Miner, i.Response.BodyContents); err != nil {
		return
	}

	// If we have a valid payload
	if len(response.Payload) > 0 {
		err = json.Unmarshal([]byte(response.Payload), &response.Quote)
	}
	return
}

// getQuote will fire the HTTP request to retrieve the fee quote
func getQuote(client *Client, miner *Miner) (result *internalResult) {
	result = &internalResult{Miner: miner}
	result.Response = httpRequest(client, http.MethodGet, defaultProtocol+miner.URL+routeFeeQuote, miner.Token, nil)
	return
}

// getQuoteRoutine will fire getQuote as part of a WaitGroup and return
// the results into a channel
func getQuoteRoutine(wg *sync.WaitGroup, client *Client, miner *Miner, resultsChannel chan *internalResult) {
	defer wg.Done()
	resultsChannel <- getQuote(client, miner)
}
