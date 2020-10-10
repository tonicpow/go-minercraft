package minercraft

import (
	"encoding/json"
	"errors"
	"net/http"
)

/*
Example submit tx(s) response from Merchant API:

{
  "payload": "{\"apiVersion\":\"1.2.3\",\"timestamp\":\"2020-09-23T09:23:02.1369987Z\",\"minerId\":\"030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e\",\"currentHighestBlockHash\":\"76beb84c5c709b6db91e2ea871c197d1ae15b47cd2caae780a26321eb2e9a290\",\"currentHighestBlockHeight\":151,\"txSecondMempoolExpiry\":0,\"txs\":[{\"txid\":\"a48f0c1c06024bb6d4d649b77e2f9b2636fa6a609ee85fc76603cecb5de0ffcb\",\"returnResult\":\"failure\",\"resultDescription\":\"Missing inputs\",\"conflictedWith\":[{\"txid\":\"ea559ed57a2f7dfdb7f33cb8601de113d782cea169eced6fb460c3e60066c323\",\"size\":191,\"hex\":\"01000000010fff3df96efcbc5637d923e27de88c45ed26392c4b791db46b8f69bd1192fc53010000006a47304402207c03d95ba831fe3686a5a081e81346312c67e33a845ed3be4fd053b3898556c0022074c24b2ae143756b5bcd20857e1d96549fd31f9b647599e255c5fcdb776fdcb04121027ae06a5b3fe1de495fa9d4e738e48810b8b06fa6c959a5305426f78f42b48f8cffffffff0198929800000000001976a91482932cf55b847ffa52832d2bbec2838f658f226788ac00000000\"}]},{\"txid\":\"65d11409d204ea80c81152d4c12ddbd37df72a0ee73828497c14cd6a0086eaf3\",\"returnResult\":\"success\"}],\"failureCount\":1}",
  "signature": "3044022076680491fa27e832a2002abce70fcccc7b859ed21c3db87c8fa0dbf143809c59022024898bbc0ef220201dd636f9fff9ee21ce90226ff327b59a505a622c9b42717b",
  "publicKey": "030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e",
  "encoding": "UTF-8",
  "mimetype": "application/json"
}
*/

// SubmitTransactionsResponse is the raw response from the Merchant API request
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi/tree/v1.2-beta#Submit-multiple-transactions
type SubmitTransactionsResponse struct {
	JSONEnvelope
	Results *SubmissionsPayload `json:"results"` // Custom field for unmarshalled payload data
}

/*
Example SubmitTransactionsResponse.Payload (unmarshalled):

{
  "apiVersion": "1.2.3",
  "timestamp": "2020-09-23T09:23:02.1369987Z",
  "minerId": "030d1fe5c1b560efe196ba40540ce9017c20daa9504c4c4cec6184fc702d9f274e",
  "currentHighestBlockHash": "76beb84c5c709b6db91e2ea871c197d1ae15b47cd2caae780a26321eb2e9a290",
  "currentHighestBlockHeight": 151,
  "txSecondMempoolExpiry": 0,
  "txs": [
    {
      "txid": "a48f0c1c06024bb6d4d649b77e2f9b2636fa6a609ee85fc76603cecb5de0ffcb",
      "returnResult": "failure",
      "resultDescription": "Missing inputs",
      "conflictedWith": [
        {
          "txid": "ea559ed57a2f7dfdb7f33cb8601de113d782cea169eced6fb460c3e60066c323",
          "size": 191,
          "hex": "01000000010fff3df96efcbc5637d923e27de88c45ed26392c4b791db46b8f69bd1192fc53010000006a47304402207c03d95ba831fe3686a5a081e81346312c67e33a845ed3be4fd053b3898556c0022074c24b2ae143756b5bcd20857e1d96549fd31f9b647599e255c5fcdb776fdcb04121027ae06a5b3fe1de495fa9d4e738e48810b8b06fa6c959a5305426f78f42b48f8cffffffff0198929800000000001976a91482932cf55b847ffa52832d2bbec2838f658f226788ac00000000"
        }
      ]
    },
    {
      "txid": "65d11409d204ea80c81152d4c12ddbd37df72a0ee73828497c14cd6a0086eaf3",
      "returnResult": "success"
    }
  ],
  "failureCount": 1
}
*/

// SubmissionsPayload is the unmarshalled version of the payload envelope
type SubmissionsPayload struct {
	APIVersion                string          `json:"apiVersion"`
	Timestamp                 string          `json:"timestamp"`
	MinerID                   string          `json:"minerId"`
	CurrentHighestBlockHash   string          `json:"currentHighestBlockHash"`
	CurrentHighestBlockHeight int64           `json:"currentHighestBlockHeight"`
	TxSecondMempoolExpiry     int64           `json:"txSecondMempoolExpiry"`
	Txs                       []*transactions `json:"txs"`
}

// transactions is the individual tx result in a multiple submission response
type transactions struct {
	TxID              string          `json:"txid"`
	ReturnResult      string          `json:"returnResult"`
	ResultDescription string          `json:"resultDescription"`
	ConflictedWith    []*conflictedTx `json:"conflictedWith,omitempty"`
}

// conflictedTx is returned if there is a conflict
type conflictedTx struct {
	Hex  string `json:"hex,omitempty"`
	Size string `json:"size,omitempty"`
	TxID string `json:"txid"`
}

// SubmitTransactions will fire a Merchant API request to submit multiple transaction
//
// This endpoint is used to send multiple raw transactions to a miner for inclusion in the
// next block that the miner creates. It returns a JSONEnvelope with a payload that contains
// the responses to the transaction submissions. The purpose of the envelope is to ensure
// strict consistency in the message content for the purpose of signing responses.
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi/tree/v1.2-beta#Submit-transaction
func (c *Client) SubmitTransactions(miner *Miner, txs []*Transaction) (*SubmitTransactionsResponse, error) {

	// Make sure we have a valid miner
	if miner == nil {
		return nil, errors.New("miner was nil")
	}

	// Make the HTTP request
	result := submitTransactions(c, miner, txs)
	if result.Response.Error != nil {
		return nil, result.Response.Error
	}

	// Parse the response
	response, err := result.parseSubmissions()
	if err != nil {
		return nil, err
	}

	// Valid query?
	if response.Results == nil {
		return nil, errors.New("failed getting submissions response from: " + miner.Name)
	}

	// Return the fully parsed response
	return &response, nil
}

// parseSubmissions will convert the HTTP response into a struct and also unmarshal the payload JSON data
func (i *internalResult) parseSubmissions() (response SubmitTransactionsResponse, err error) {

	// Process the initial response payload
	if err = response.process(i.Miner, i.Response.BodyContents); err != nil {
		return
	}

	// If we have a valid payload
	if len(response.Payload) > 0 {
		err = json.Unmarshal([]byte(response.Payload), &response.Results)
	}
	return
}

// submitTransactions will fire the HTTP request to submit multiple transactions
func submitTransactions(client *Client, miner *Miner, txs []*Transaction) (result *internalResult) {
	result = &internalResult{Miner: miner}
	data, _ := json.Marshal(txs) // Ignoring error - if it fails, the submission would also fail
	result.Response = httpRequest(client, http.MethodPost, defaultProtocol+miner.URL+routeSubmitTxs, miner.Token, data)
	return
}
