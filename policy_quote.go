package minercraft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/libsv/go-bt/v2"
	"github.com/tonicpow/go-minercraft/v2/apis/arc"
	"github.com/tonicpow/go-minercraft/v2/apis/mapi"
)

// PolicyQuoteModelAdapter is the interface for the adapter to get the policy quote response
type PolicyQuoteModelAdapter interface {
	GetPolicyData() *PolicyPayload
}

// PolicyQuoteMapiAdapter is the adapter for the mAPI response
type PolicyQuoteMapiAdapter struct {
	*mapi.PolicyQuoteModel
}

// PolicyQuoteArcAdapter is the adapter for the Arc response
type PolicyQuoteArcAdapter struct {
	*arc.PolicyQuoteModel
}

// UnifiedPolicy is the unmarshalled version of the policy
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

// UnifiedFeePayload is the unmarshalled version of the payload envelope
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

	route, err := ActionRouteByAPIType(PolicyQuote, c.apiType)
	if err != nil {
		return nil, err
	}

	// Make the HTTP request
	result := getQuote(ctx, c, miner, route)
	if result.Response.Error != nil {
		return nil, result.Response.Error
	}

	quoteResponse := &PolicyQuoteResponse{
		JSONEnvelope: JSONEnvelope{
			APIType: c.apiType,
			Miner:   result.Miner,
		},
	}

	var modelAdapter PolicyQuoteModelAdapter

	switch c.apiType {
	case MAPI:
		model := &mapi.PolicyQuoteModel{}
		err = quoteResponse.process(result.Miner, result.Response.BodyContents)
		if err != nil || len(quoteResponse.Payload) <= 0 {
			return nil, err
		}

		err = json.Unmarshal([]byte(quoteResponse.Payload), model)
		if err != nil {
			return nil, err
		}

		if model.Fees == nil || len(model.Fees) <= 0 {
			return nil, errors.New("empty fees")
		}

		modelAdapter = &PolicyQuoteMapiAdapter{PolicyQuoteModel: model}
	case Arc:
		model := &arc.PolicyQuoteModel{}
		err = json.Unmarshal(result.Response.BodyContents, model)
		if err != nil {
			return nil, err
		}

		modelAdapter = &PolicyQuoteArcAdapter{PolicyQuoteModel: model}
	default:
		return nil, fmt.Errorf("unknown API type: %s", c.apiType)
	}

	quoteResponse.Quote = modelAdapter.GetPolicyData()

	// Valid?
	if quoteResponse.Quote == nil {
		return nil, errors.New("failed getting quote response from: " + miner.Name)
	}

	isValid, err := quoteResponse.IsValid()
	if err != nil {
		return nil, err
	}

	quoteResponse.Validated = isValid

	// Return the fully parsed response
	return quoteResponse, nil
}

// GetPolicyData will return the policy data from the mapi adapter
func (a *PolicyQuoteMapiAdapter) GetPolicyData() *PolicyPayload {
	// Creates instance of UnifiedFeePayload
	feePayload := UnifiedFeePayload{
		FeePayloadFields: mapi.FeePayloadFields{
			APIVersion:                a.APIVersion,
			Timestamp:                 a.Timestamp,
			ExpirationTime:            a.ExpiryTime,
			MinerID:                   a.MinerID,
			CurrentHighestBlockHash:   a.CurrentHighestBlockHash,
			CurrentHighestBlockHeight: a.CurrentHighestBlockHeight,
			MinerReputation:           nil,
		},
		Fees: make([]*bt.Fee, len(a.Fees)),
	}

	for i, mapiFee := range a.Fees {
		feePayload.Fees[i] = &bt.Fee{
			FeeType:   bt.FeeType(mapiFee.FeeType),
			MiningFee: bt.FeeUnit(mapiFee.MiningFee),
			RelayFee:  bt.FeeUnit(mapiFee.RelayFee),
		}
	}

	callbacks := make([]*mapi.PolicyCallback, len(a.Callbacks))
	for i, cb := range a.Callbacks {
		//nolint:gosec,exportloopref // ignore those linter errors
		callbacks[i] = &cb
	}

	policyPayload := &PolicyPayload{
		UnifiedFeePayload: feePayload,
		Callbacks:         callbacks,
		Policies: &UnifiedPolicy{
			AcceptNonStdOutputs:             a.Policies.AcceptNonStdOutputs,
			DataCarrier:                     a.Policies.DataCarrier,
			DataCarrierSize:                 a.Policies.DataCarrierSize,
			LimitAncestorCount:              a.Policies.LimitAncestorCount,
			LimitCpfpGroupMembersCount:      a.Policies.LimitCpfpGroupMembersCount,
			MaxNonStdTxValidationDuration:   a.Policies.MaxNonStdTxValidationDuration,
			MaxScriptNumLengthPolicy:        a.Policies.MaxScriptNumLengthPolicy,
			MaxScriptSizePolicy:             a.Policies.MaxScriptSizePolicy,
			MaxStackMemoryUsagePolicy:       a.Policies.MaxStackMemoryUsagePolicy,
			MaxStdTxValidationDuration:      a.Policies.MaxStdTxValidationDuration,
			MaxTxSizePolicy:                 a.Policies.MaxTxSizePolicy,
			SkipScriptFlags:                 a.Policies.SkipScriptFlags,
			MaxConsolidationFactor:          a.Policies.MaxConsolidationFactor,
			MaxConsolidationInputScriptSize: a.Policies.MaxConsolidationInputScriptSize,
			MinConfConsolidationInput:       a.Policies.MinConfConsolidationInput,
			AcceptNonStdConsolidationInput:  a.Policies.AcceptNonStdConsolidationInput,
			MaxTxSigOpsCount:                0,
		},
	}

	return policyPayload
}

// GetPolicyData will return the policy data from the arc adapter
func (a *PolicyQuoteArcAdapter) GetPolicyData() *PolicyPayload {

	feePayload := UnifiedFeePayload{
		FeePayloadFields: mapi.FeePayloadFields{
			APIVersion:                "",
			Timestamp:                 a.Timestamp,
			ExpirationTime:            "",
			MinerID:                   "",
			CurrentHighestBlockHash:   "",
			CurrentHighestBlockHeight: 0,
			MinerReputation:           nil,
		},
		Fees: []*bt.Fee{a.Policy.MiningFee},
	}

	policyPayload := &PolicyPayload{
		UnifiedFeePayload: feePayload,
		Callbacks:         nil,
		Policies: &UnifiedPolicy{
			AcceptNonStdOutputs:             false,
			DataCarrier:                     false,
			DataCarrierSize:                 0,
			LimitAncestorCount:              0,
			LimitCpfpGroupMembersCount:      0,
			MaxNonStdTxValidationDuration:   0,
			MaxScriptNumLengthPolicy:        0,
			MaxScriptSizePolicy:             a.Policy.MaxScriptSizePolicy,
			MaxStackMemoryUsagePolicy:       0,
			MaxStdTxValidationDuration:      0,
			MaxTxSizePolicy:                 a.Policy.MaxTxSizePolicy,
			SkipScriptFlags:                 nil,
			MaxConsolidationFactor:          0,
			MaxConsolidationInputScriptSize: 0,
			MinConfConsolidationInput:       0,
			AcceptNonStdConsolidationInput:  false,
			MaxTxSigOpsCount:                a.Policy.MaxTxSigOpsCount,
		},
	}

	return policyPayload
}
