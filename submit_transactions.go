package minercraft

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Reference: https://github.com/bitcoin-sv-specs/brfc-merchantapi#4-submit-multiple-transactions

type (

	// RawSubmitTransactionsResponse is the response returned from mapi where payload is a string.
	RawSubmitTransactionsResponse struct {
		TxsPayload string `json:"payload"`
		Signature  string `json:"signature"`
		Publickey  string `json:"publicKey"`
		Encoding   string `json:"encoding"`
		Mimetype   string `json:"mimetype"`
	}

	// SubmitTransactionsResponse is the formatted response which converts payload string to payloads.
	SubmitTransactionsResponse struct {
		TxsPayload TxsPayload `json:"payload"`
		Signature  string     `json:"signature"`
		Publickey  string     `json:"publicKey"`
		Encoding   string     `json:"encoding"`
		Mimetype   string     `json:"mimetype"`
	}

	// TxsPayload is the structure of the json payload string in the MapiResponse.
	TxsPayload struct {
		Apiversion                string    `json:"apiVersion"`
		Timestamp                 time.Time `json:"timestamp"`
		Minerid                   string    `json:"minerId"`
		Currenthighestblockhash   string    `json:"currentHighestBlockHash"`
		Currenthighestblockheight int       `json:"currentHighestBlockHeight"`
		Txsecondmempoolexpiry     int       `json:"txSecondMempoolExpiry"`
		Txs                       []Tx      `json:"txs"`
		Failurecount              int       `json:"failureCount"`
	}

	// Tx is the transaction format in the mapi txs response.
	Tx struct {
		Txid              string           `json:"txid"`
		Returnresult      string           `json:"returnResult"`
		Resultdescription string           `json:"resultDescription"`
		Conflictedwith    []ConflictedWith `json:"conflictedWith,omitempty"`
	}
)

func (c *Client) SubmitTransactions(ctx context.Context, miner *Miner, txs *[]Transaction) (*SubmitTransactionsResponse, error) {
	if miner == nil {
		return nil, errors.New("miner was nil")
	}

	if len(*txs) <= 0 {
		return nil, errors.New("no transactions")
	}

	data, err := json.Marshal(txs)
	if err != nil {
		return nil, err
	}

	response := httpRequest(ctx, c, &httpPayload{
		Method: http.MethodPost,
		URL:    miner.URL + routeSubmitTxs,
		Token:  miner.Token,
		Data:   data,
	})

	if response.Error != nil {
		return nil, err
	}

	var raw *RawSubmitTransactionsResponse
	err = json.Unmarshal(response.BodyContents, &raw)
	if err != nil {
		return nil, err
	}

	result := &SubmitTransactionsResponse{
		Signature: raw.Signature,
		Publickey: raw.Publickey,
		Encoding:  raw.Encoding,
		Mimetype:  raw.Mimetype,
	}
	err = json.Unmarshal([]byte(raw.TxsPayload), &result.TxsPayload)
	if err != nil {
		return nil, err
	}

	return result, err
}
