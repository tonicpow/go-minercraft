package minercraft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tonicpow/go-minercraft/apis/arc"
	"github.com/tonicpow/go-minercraft/apis/mapi"
)

const (
	// MerkleFormatTSC can be set when calling SubmitTransaction to request a MerkleProof in TSC format.
	MerkleFormatTSC = "TSC"
)

type SubmitTxModelAdapter interface {
	GetSubmitTxResponse() *UnifiedSubmissionPayload
}

type SubmitTxMapiAdapter struct {
	*mapi.SubmitTxModel
}

type SubmitTxArcAdapter struct {
	*arc.SubmitTxModel
}

// SubmitTransactionResponse is the raw response from the Merchant API request
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi#3-submit-transaction
type SubmitTransactionResponse struct {
	JSONEnvelope
	Results *UnifiedSubmissionPayload `json:"results"` // Custom field for unmarshalled payload data
}

/*
Example SubmitTransactionResponse.Payload (unmarshalled):

{
  "apiVersion": "1.2.3",
  "conflictedWith": ""
  "currentHighestBlockHash": "71a7374389afaec80fcabbbf08dcd82d392cf68c9a13fe29da1a0c853facef01",
  "currentHighestBlockHeight": 207,
  "minerId": "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031",
  "resultDescription": "",
  "returnResult": "success",
  "timestamp": "2020-01-15T11:40:29.826Z",
  "txid": "6bdbcfab0526d30e8d68279f79dff61fb4026ace8b7b32789af016336e54f2f0",
  "txSecondMempoolExpiry": 0,
}
*/

// UnifiedSubmissionPayload is the unmarshalled version of the payload envelope
type UnifiedSubmissionPayload struct {
	// mAPI
	APIVersion                string                 `json:"apiVersion"`
	ConflictedWith            []*mapi.ConflictedWith `json:"conflictedWith"`
	CurrentHighestBlockHash   string                 `json:"currentHighestBlockHash"`
	CurrentHighestBlockHeight int64                  `json:"currentHighestBlockHeight"`
	MinerID                   string                 `json:"minerId"`
	ResultDescription         string                 `json:"resultDescription"`
	ReturnResult              string                 `json:"returnResult"`
	Timestamp                 string                 `json:"timestamp"`
	TxID                      string                 `json:"txid"`
	TxSecondMempoolExpiry     int64                  `json:"txSecondMempoolExpiry"`
	// FailureRetryable if true indicates the tx can be resubmitted to mAPI.
	FailureRetryable bool `json:"failureRetryable"`

	// Arc
	BlockHash   string       `json:"blockHash,omitempty"`
	BlockHeight int64        `json:"blockHeight,omitempty"`
	ExtraInfo   string       `json:"extraInfo,omitempty"`
	Status      int          `json:"status,omitempty"`
	Title       string       `json:"title,omitempty"`
	TxStatus    arc.TxStatus `json:"txStatus,omitempty"`
}

// Transaction is the body contents in the "submit transaction" request
type Transaction struct {
	CallBackEncryption string       `json:"callBackEncryption,omitempty"`
	CallBackToken      string       `json:"callBackToken,omitempty"`
	CallBackURL        string       `json:"callBackUrl,omitempty"`
	DsCheck            bool         `json:"dsCheck,omitempty"`
	MerkleFormat       string       `json:"merkleFormat,omitempty"`
	MerkleProof        bool         `json:"merkleProof,omitempty"`
	RawTx              string       `json:"rawtx"`
	WaitForStatus      arc.TxStatus `json:"waitForStatus,omitempty"`
}

// SubmitTransaction will fire a Merchant API request to submit a given transaction
//
// This endpoint is used to send a raw transaction to a miner for inclusion in the next block
// that the miner creates. It returns a JSONEnvelope with a payload that contains the response to the
// transaction submission. The purpose of the envelope is to ensure strict consistency in the
// message content for the purpose of signing responses.
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi#3-submit-transaction
func (c *Client) SubmitTransaction(ctx context.Context, miner *Miner, tx *Transaction) (*SubmitTransactionResponse, error) {

	// Make sure we have a valid miner
	if miner == nil {
		return nil, errors.New("miner was nil")
	}

	// Make the HTTP request
	result, err := submitTransaction(ctx, c, miner, tx)
	if err != nil {
		return nil, err
	}
	if result.Response.Error != nil {
		return nil, result.Response.Error
	}

	submitResponse := &SubmitTransactionResponse{
		JSONEnvelope: JSONEnvelope{
			ApiType: c.apiType,
			Miner:   result.Miner,
		},
	}

	var modelAdapter SubmitTxModelAdapter

	switch c.apiType {
	case MAPI:
		model := &SubmitTxMapiAdapter{}
		err := submitResponse.process(result.Miner, result.Response.BodyContents)
		if err != nil || len(submitResponse.Payload) <= 0 {
			return nil, err
		}

		err = json.Unmarshal([]byte(submitResponse.Payload), model)
		if err != nil {
			return nil, err
		}

		modelAdapter = &SubmitTxMapiAdapter{SubmitTxModel: model.SubmitTxModel}
	case Arc:
		model := &SubmitTxArcAdapter{}
		err := json.Unmarshal(result.Response.BodyContents, model)
		if err != nil {
			return nil, err
		}

		modelAdapter = &SubmitTxArcAdapter{SubmitTxModel: model.SubmitTxModel}
	default:
		return nil, fmt.Errorf("unknown API type: %s", c.apiType)
	}

	submitResponse.Results = modelAdapter.GetSubmitTxResponse()

	// Valid?
	if submitResponse.Results == nil {
		return nil, errors.New("failed getting submit response from: " + miner.Name)
	}

	isValid, err := submitResponse.IsValid()
	if err != nil {
		return nil, err
	}

	submitResponse.Validated = isValid

	// Return the fully parsed response
	return submitResponse, nil
}

// submitTransaction will fire the HTTP request to submit a transaction
func submitTransaction(ctx context.Context, client *Client, miner *Miner, tx *Transaction) (*internalResult, error) {
	result := &internalResult{Miner: miner}

	api, err := MinerAPIByMinerID(client.minerAPIs, miner.MinerID, client.apiType)
	if err != nil {
		result.Response = &RequestResponse{Error: err}
		return nil, err
	}

	route, err := ActionRouteByAPIType(SubmitTx, client.apiType)
	if err != nil {
		result.Response = &RequestResponse{Error: err}
		return nil, err
	}

	submitURL := api.URL + route
	httpPayload := &httpPayload{
		Method:  http.MethodPost,
		URL:     submitURL,
		Token:   api.Token,
		Headers: make(map[string]string),
	}

	switch client.apiType {
	case MAPI:
		err = proceedMapi(tx, httpPayload)
		if err != nil {
			return nil, err
		}

	case Arc:
		err = proceedArc(tx, httpPayload)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown API type: %s", client.apiType)
	}

	result.Response = httpRequest(ctx, client, httpPayload)
	return result, nil
}

func proceedArc(tx *Transaction, httpPayload *httpPayload) error {
	body := map[string]string{
		"rawTx": tx.RawTx,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshall JSON when submitting transaction: %w", err)
	}
	httpPayload.Data = data

	if tx.MerkleProof {
		httpPayload.Headers["X-MerkleProof"] = "true"
	}

	if tx.CallBackURL != "" {
		httpPayload.Headers["X-CallbackUrl"] = tx.CallBackURL
	}

	if tx.CallBackToken != "" {
		httpPayload.Headers["X-CallbackToken"] = tx.CallBackToken
	}

	if statusCode, ok := arc.MapTxStatusToInt(tx.WaitForStatus); ok {
		httpPayload.Headers["X-WaitForStatus"] = strconv.Itoa(statusCode)
	}

	return nil
}

func proceedMapi(tx *Transaction, httpPayload *httpPayload) error {
	data, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	httpPayload.Data = data
	return nil
}

func (a *SubmitTxMapiAdapter) GetSubmitTxResponse() *UnifiedSubmissionPayload {
	return &UnifiedSubmissionPayload{
		APIVersion:                a.APIVersion,
		ConflictedWith:            a.ConflictedWith,
		CurrentHighestBlockHash:   a.CurrentHighestBlockHash,
		CurrentHighestBlockHeight: a.CurrentHighestBlockHeight,
		MinerID:                   a.MinerID,
		ResultDescription:         a.ResultDescription,
		ReturnResult:              a.ReturnResult,
		Timestamp:                 a.Timestamp,
		TxID:                      a.TxID,
		TxSecondMempoolExpiry:     a.TxSecondMempoolExpiry,
		FailureRetryable:          a.FailureRetryable,
	}
}

func (a *SubmitTxArcAdapter) GetSubmitTxResponse() *UnifiedSubmissionPayload {
	return &UnifiedSubmissionPayload{
		BlockHash:   a.BlockHash,
		BlockHeight: a.BlockHeight,
		ExtraInfo:   a.ExtraInfo,
		Status:      a.Status,
		Title:       a.Title,
		TxStatus:    a.TxStatus,
		TxID:        a.TxID,
	}
}
