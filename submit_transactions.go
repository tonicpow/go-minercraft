package minercraft

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tonicpow/go-minercraft/apis/arc"
	"github.com/tonicpow/go-minercraft/apis/mapi"
)

// Reference: https://github.com/bitcoin-sv-specs/brfc-merchantapi#5-submit-multiple-transactions

type (

	// RawSubmitTransactionsResponse is the response returned from mapi where payload is a string.
	RawSubmitTransactionsResponse struct {
		Encoding  string `json:"encoding"`
		MimeType  string `json:"mimetype"`
		Payload   string `json:"payload"`
		PublicKey string `json:"publicKey"`
		Signature string `json:"signature"`
	}

	// SubmitTransactionsResponse is the formatted response which converts payload string to payloads.
	SubmitTransactionsResponse struct {
		Encoding  string            `json:"encoding"`
		MimeType  string            `json:"mimetype"`
		Payload   UnifiedTxsPayload `json:"payload"`
		PublicKey string            `json:"publicKey"`
		Signature string            `json:"signature"`
	}

	// TxsPayload is the structure of the json payload string in the MapiResponse.
	UnifiedTxsPayload struct {
		APIVersion                string    `json:"apiVersion"`
		CurrentHighestBlockHash   string    `json:"currentHighestBlockHash"`
		CurrentHighestBlockHeight int       `json:"currentHighestBlockHeight"`
		FailureCount              int       `json:"failureCount"`
		MinerID                   string    `json:"minerId"`
		Timestamp                 time.Time `json:"timestamp"`
		Txs                       []Tx      `json:"txs"`
		TxSecondMempoolExpiry     int       `json:"txSecondMempoolExpiry"`

		// Arc
		BlockHash   string       `json:"blockHash,omitempty"`
		BlockHeight int64        `json:"blockHeight,omitempty"`
		ExtraInfo   string       `json:"extraInfo,omitempty"`
		Status      int          `json:"status,omitempty"`
		Title       string       `json:"title,omitempty"`
		TxStatus    arc.TxStatus `json:"txStatus,omitempty"`
	}

	// Tx is the transaction format in the mapi txs response.
	Tx struct {
		ConflictedWith    []mapi.ConflictedWith `json:"conflictedWith,omitempty"`
		ResultDescription string                `json:"resultDescription"`
		ReturnResult      string                `json:"returnResult"`
		TxID              string                `json:"txid"`
		// FailureRetryable if true indicates the tx can be resubmitted to mAPI.
		FailureRetryable bool `json:"failureRetryable"`
	}
)

// SubmitTransactions is used for submitting batched transactions
//
// Reference: https://github.com/bitcoin-sv-specs/brfc-merchantapi#5-submit-multiple-transactions
func (c *Client) SubmitTransactions(ctx context.Context, miner *Miner, txs []Transaction) (*SubmitTransactionsResponse, error) {
	if miner == nil {
		return nil, errors.New("miner was nil")
	}

	if len(txs) <= 0 {
		return nil, errors.New("no transactions")
	}

	if response.Error != nil {
		return nil, response.Error
	}

	var raw RawSubmitTransactionsResponse
	if err = json.Unmarshal(response.BodyContents, &raw); err != nil {
		return nil, err
	}

	result := &SubmitTransactionsResponse{
		Signature: raw.Signature,
		PublicKey: raw.PublicKey,
		Encoding:  raw.Encoding,
		MimeType:  raw.MimeType,
	}

	if err = json.Unmarshal([]byte(raw.Payload), &result.Payload); err != nil {
		return nil, err
	}

	return result, err
}

func submitTransactions(ctx context.Context, client *Client, miner *Miner, txs []Transaction) (*UnifiedTxsPayload, error) {

	result := UnifiedTxsPayload{}
	api, err := MinerAPIByMinerID(client.minerAPIs, miner.MinerID, client.apiType)
	if err != nil {
		return nil, err
	}

	route, err := ActionRouteByAPIType(SubmitTx, client.apiType)
	if err != nil {
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
		err = proceedMapiSubmitTxs(txs, httpPayload)
		if err != nil {
			return nil, err
		}

	case Arc:
		err = proceedArcSubmitTxs(txs, httpPayload)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown API type: %s", client.apiType)
	}

	response := httpRequest(ctx, client, httpPayload)
	return result, nil
}

func proceedArcSubmitTxs(txs []Transaction, httpPayload *httpPayload) error {
	var rawTxs []string
	for _, tx := range txs {
		rawTxs = append(rawTxs, tx.RawTx)
	}

	body := map[string]interface{}{
		"rawTx": rawTxs,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON when submitting transactions: %w", err)
	}

	httpPayload.Data = data

	if txs[0].MerkleProof {
		httpPayload.Headers["X-MerkleProof"] = "true"
	}

	if txs[0].CallBackURL != "" {
		httpPayload.Headers["X-CallbackUrl"] = txs[0].CallBackURL
	}

	if txs[0].CallBackToken != "" {
		httpPayload.Headers["X-CallbackToken"] = txs[0].CallBackToken
	}

	if statusCode, ok := arc.MapTxStatusToInt(txs[0].WaitForStatus); ok {
		httpPayload.Headers["X-WaitForStatus"] = strconv.Itoa(statusCode)
	}

	return nil
}

func proceedMapiSubmitTxs(txs []Transaction, httpPayload *httpPayload) error {
	data, err := json.Marshal(txs)
	if err != nil {
		return err
	}

	httpPayload.Data = data
	return nil
}
