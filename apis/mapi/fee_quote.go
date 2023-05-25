package mapi

import (
	"fmt"
	"strings"

	"github.com/libsv/go-bt/v2"
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
	"payload": "{\"apiVersion\":\"1.4.0\",\"timestamp\":\"2020-10-07T21:13:04.335Z\",\"expiryTime\":\"2020-10-07T21:23:04.335Z\",\"minerId\":\"0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087\",\"currentHighestBlockHash\":\"000000000000000000edb30c3bbbc8e6a07e522e85522e6a213f7e933e6e2d8d\",\"currentHighestBlockHeight\":655874,\"minerReputation\":null,\"fees\":[{\"feeType\":\"standard\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}},{\"feeType\":\"data\",\"miningFee\":{\"satoshis\":500,\"bytes\":1000},\"relayFee\":{\"satoshis\":250,\"bytes\":1000}}]}",
	"signature": "304402206443bea5bdd98a16e23eb61c36b4b998bd68ceb9c84983c7e695e267b21a30440220191571e9b9632c8337d9196723ca20eefa63966ef6360170db0e57a04047453f",
	"publicKey": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087",
	"encoding": "UTF-8",
	"mimetype": "application/json"
}
*/

/*
Example FeeQuoteResponse.Payload (unmarshalled):

{
  "apiVersion": "1.4.0",
  "timestamp": "2020-10-07T21:13:04.335Z",
  "expiryTime": "2020-10-07T21:23:04.335Z",
  "minerId": "0211ccfc29e3058b770f3cf3eb34b0b2fd2293057a994d4d275121be4151cdf087",
  "currentHighestBlockHash": "000000000000000000edb30c3bbbc8e6a07e522e85522e6a213f7e933e6e2d8d",
  "currentHighestBlockHeight": 655874,
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
	FeePayloadFields
	Fees []*bt.Fee `json:"fees"`
}

type (

	// rawFeePayload is the unmarshalled version of the payload envelope
	RawFeePayload struct {
		FeePayloadFields
		Callbacks []*PolicyCallback `json:"callbacks"` // IP addresses of double-spend notification servers such as mAPI reference implementation
		Fees      []*feeObj         `json:"fees"`
	}

	// feePayloadFields are the same fields in both payloads
	FeePayloadFields struct {
		APIVersion                string      `json:"apiVersion"`
		Timestamp                 string      `json:"timestamp"`
		ExpirationTime            string      `json:"expiryTime"`
		MinerID                   string      `json:"minerId"`
		CurrentHighestBlockHash   string      `json:"currentHighestBlockHash"`
		CurrentHighestBlockHeight uint64      `json:"currentHighestBlockHeight"`
		MinerReputation           interface{} `json:"minerReputation"` // Not sure what this value is
	}

	// feeUnit displays the amount of Satoshis needed
	// for a specific amount of Bytes in a transaction
	// see https://github.com/bitcoin-sv-specs/brfc-merchantapi#expanded-payload-1
	FeeUnit struct {
		Satoshis int `json:"satoshis"` // Fee in satoshis of the amount of Bytes
		Bytes    int `json:"bytes"`    // Number of bytes that the Fee covers
	}

	// feeObj displays the MiningFee as well as the RelayFee for a specific
	// FeeType, for example 'standard' or 'data'
	// see https://github.com/bitcoin-sv-specs/brfc-merchantapi#expanded-payload-1
	feeObj struct {
		FeeType   string  `json:"feeType"` // standard || data
		MiningFee FeeUnit `json:"miningFee"`
		RelayFee  FeeUnit `json:"relayFee"` // Fee for retaining Tx in secondary mempool
	}
)

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
		if string(fee.FeeType) != feeType {
			continue
		}

		// Multiply & Divide
		var calcFee uint64
		if strings.EqualFold(feeCategory, FeeCategoryMining) {
			calcFee = (uint64(fee.MiningFee.Satoshis) * txBytes) / uint64(fee.MiningFee.Bytes)
		} else {
			calcFee = (uint64(fee.RelayFee.Satoshis) * txBytes) / uint64(fee.RelayFee.Bytes)
		}

		// Check for zero
		if calcFee != 0 {
			return calcFee, nil
		}

		// If txBytes is zero this error will occur
		return 1, fmt.Errorf("warning: fee calculation was 0")
	}

	// No fee type found in the slice of fees
	return 1, fmt.Errorf("feeType %s is not found in fees", feeType)
}

// GetFee will return the fee associated to the type (standard, data)
func (f *FeePayload) GetFee(feeType string) *bt.Fee {

	// Loop the fees for the given type
	for index, fee := range f.Fees {
		if string(fee.FeeType) == feeType {
			return f.Fees[index]
		}
	}

	return nil
}
