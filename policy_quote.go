package minercraft

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/libsv/go-bt/v2"
	"github.com/tonicpow/go-minercraft/apis/arc"
	"github.com/tonicpow/go-minercraft/apis/mapi"
)

type PolicyQuoteModelAdapter interface {
	GetPolicyData() *PolicyPayload
}

type PolicyQuoteMapiAdapter struct {
	*mapi.PolicyQuoteModel
}

type PolicyQuoteArcAdapter struct {
	*arc.PolicyQuoteModel
}

type UnifiedPolicy struct {
	AcceptNonStdOutputs             bool              `json:"acceptnonstdoutputs"`
	DataCarrier                     bool              `json:"datacarrier"`
	DataCarrierSize                 uint32            `json:"datacarriersize"`
	LimitAncestorCount              uint32            `json:"limitancestorcount"`
	LimitCpfpGroupMembersCount      uint32            `json:"limitcpfpgroupmemberscount"`
	MaxNonStdTxValidationDuration   uint32            `json:"maxnonstdtxvalidationduration"`
	MaxScriptNumLengthPolicy        uint32            `json:"maxscriptnumlengthpolicy"`
	MaxScriptSizePolicy             uint32            `json:"maxscriptsizepolicy"`
	MaxStackMemoryUsagePolicy       uint64            `json:"maxstackmemoryusagepolicy"`
	MaxStdTxValidationDuration      uint32            `json:"maxstdtxvalidationduration"`
	MaxTxSizePolicy                 uint32            `json:"maxtxsizepolicy"`
	SkipScriptFlags                 []mapi.ScriptFlag `json:"skipscriptflags"`
	MaxConsolidationFactor          uint32            `json:"minconsolidationfactor"`
	MaxConsolidationInputScriptSize uint32            `json:"maxconsolidationinputscriptsize"`
	MinConfConsolidationInput       uint32            `json:"minconfconsolidationinput"`
	AcceptNonStdConsolidationInput  bool              `json:"acceptnonstdconsolidationinput"`

	// Additional fields for Policy in API2
	MaxTxSigOpsCount uint32 `json:"maxtxsigopscount"`
}

type UnifiedFeePayload struct {
	mapi.FeePayloadFields
	Fees []*bt.Fee `json:"fees"`
}

// PolicyPayload is the unmarshalled version of the payload envelope
type PolicyPayload struct {
	UnifiedFeePayload                        // Inherit the same structure as the fee payload
	Callbacks         []*mapi.PolicyCallback `json:"callbacks"` // IP addresses of double-spend notification servers such as mAPI reference implementation
	Policies          *UnifiedPolicy         `json:"policies"`  // values of miner policies as configured by the mAPI reference implementation administrator
}

// PolicyQuoteResponse is the raw response from the API request
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi#1-get-policy-quote
type PolicyQuoteResponse struct {
	JSONEnvelope
	Quote *PolicyPayload `json:"quote"` // Custom field for unmarshalled payload data
}

// PolicyQuote will fire a Merchant&Arc API request to retrieve the policy from a given miner
//
// This endpoint is used to get the different policies quoted by a miner.
// It returns a JSONEnvelope with a payload that contains the policies used by a specific BSV miner.
// The purpose of the envelope is to ensure strict consistency in the message content for
// the purpose of signing responses. This is a superset of the fee quote service, as it also
// includes information on DSNT IP addresses and miner policies.
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi#1-get-policy-quote
// Specs: https://docs.gorillapool.io/arc/api.html#get-the-policy-settings
func (c *Client) PolicyQuote(ctx context.Context, miner *Miner) (*PolicyQuoteResponse, error) {

	// Make sure we have a valid miner
	if miner == nil {
		return nil, errors.New("miner was nil")
	}

	// Make the HTTP request
	result := getQuote(ctx, c, miner, mAPIRoutePolicyQuote)
	if result.Response.Error != nil {
		return nil, result.Response.Error
	}

	// Parse the response
	response, err := result.parsePolicyQuote()
	if err != nil {
		return nil, err
	}

	// Valid?
	if response.Quote == nil || len(response.Quote.Fees) == 0 {
		return nil, errors.New("failed getting policy from: " + miner.Name)
	}

	// Return the fully parsed response
	return &response, nil
}

// parsePolicyQuote will convert the HTTP response into a struct and also unmarshal the payload JSON data
func (i *internalResult) parsePolicyQuote() (response PolicyQuoteResponse, err error) {

	// Process the initial response payload
	if err = response.process(i.Miner, i.Response.BodyContents); err != nil {
		return
	}

	// If we have a valid payload
	if len(response.Payload) > 0 {
		if err = json.Unmarshal([]byte(response.Payload), &response.Quote); err != nil {
			return
		}
		if response.Quote != nil &&
			len(response.Quote.Fees) > 0 &&
			len(response.Quote.Fees[0].FeeType) == 0 { // This is an issue because go-bt json field is stripping the types
			response.Quote.Fees[0].FeeType = mapi.FeeTypeStandard
			response.Quote.Fees[1].FeeType = mapi.FeeTypeData
		}
	}
	return
}
