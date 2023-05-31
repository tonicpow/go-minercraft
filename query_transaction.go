package minercraft

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/libsv/go-bc"
	"github.com/tonicpow/go-minercraft/apis/arc"
	"github.com/tonicpow/go-minercraft/apis/mapi"
)

// QueryTransactionSuccess is on success
const QueryTransactionSuccess = "success"

// QueryTransactionFailure is on failure
const QueryTransactionFailure = "failure"

// QueryTransactionInMempoolFailure in mempool but not in a block yet
const QueryTransactionInMempoolFailure = "Transaction in mempool but not yet in block"

type QueryTxModelAdapter interface {
	GetQueryTxResponse() *QueryTxResponse
}

type QueryTxMapiAdapter struct {
	*mapi.QueryTxModel
}

type QueryTxArcAdapter struct {
	*arc.QueryTxModel
}

/*
Example query tx response from Merchant API:

{
  "payload": "{\"apiVersion\":\"1.2.3\",\"timestamp\":\"2020-01-15T11:41:29.032Z\",\"returnResult\":\"failure\",\"resultDescription\":\"Transaction in mempool but not yet in block\",\"blockHash\":\"\",\"blockHeight\":0,\"minerId\":\"03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031\",\"confirmations\":0,\"txSecondMempoolExpiry\":0}",
  "signature": "3045022100f78a6ac49ef38fbe68db609ff194d22932d865d93a98ee04d2ecef5016872ba50220387bf7e4df323bf4a977dd22a34ea3ad42de1a2ec4e5af59baa13258f64fe0e5",
  "publicKey": "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031",
  "encoding": "UTF-8",
  "mimetype": "application/json"
}
*/

// QueryTransactionResponse is the raw response from the Merchant API request
//
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi#4-query-transaction-status
type QueryTransactionResponse struct {
	JSONEnvelope
	Query *QueryTxResponse `json:"query"` // Custom field for unmarshalled payload data
}

/*
Example QueryTransactionResponse.Payload (unmarshalled):

Failure - in mempool but not in block
{
  "apiVersion": "1.2.3",
  "timestamp": "2020-01-15T11:41:29.032Z",
  "txid": "6bdbcfab0526d30e8d68279f79dff61fb4026ace8b7b32789af016336e54f2f0",
  "returnResult": "failure",
  "resultDescription": "Transaction in mempool but not yet in block",
  "blockHash": "",
  "blockHeight": 0,
  "minerId": "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031",
  "confirmations": 0,
  "txSecondMempoolExpiry": 0
}

Success - added to block
{
  "apiVersion": "1.2.3",
  "timestamp": "2020-01-15T12:09:37.394Z",
  "txid": "6bdbcfab0526d30e8d68279f79dff61fb4026ace8b7b32789af016336e54f2f0",
  "returnResult": "success",
  "resultDescription": "",
  "blockHash": "745093bb0c80780092d4ce6926e0caa753fe3accdc09c761aee89bafa85f05f4",
  "blockHeight": 208,
  "minerId": "03fcfcfcd0841b0a6ed2057fa8ed404788de47ceb3390c53e79c4ecd1e05819031",
  "confirmations": 2,
  "txSecondMempoolExpiry": 0
}
*/

// QueryTxResponse is the unmarshalled version of the payload envelope
type QueryTxResponse struct {
	// Pola wspólne dla obu typów API
	BlockHash   string `json:"blockHash,omitempty"`
	BlockHeight int64  `json:"blockHeight,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
	TxID        string `json:"txid,omitempty"`
	// ArcAPI specific fields
	TxStatus arc.TxStatus `json:"txStatus,omitempty"`
	// mAPI specific fields
	APIVersion            string          `json:"apiVersion,omitempty"`
	ReturnResult          string          `json:"returnResult,omitempty"`
	ResultDescription     string          `json:"resultDescription,omitempty"`
	MinerID               string          `json:"minerId,omitempty"`
	Confirmations         int64           `json:"confirmations,omitempty"`
	TxSecondMempoolExpiry int64           `json:"txSecondMempoolExpiry,omitempty"`
	MerkleProof           *bc.MerkleProof `json:"merkleProof,omitempty"`
}

// QueryTransactionOptFunc defines an optional argument that can be passed to the QueryTransaction method.
type QueryTransactionOptFunc func(o *queryTransactionOpts)

type queryTransactionOpts struct {
	includeProof bool
	merkleFormat string
}

func defaultQueryOpts() *queryTransactionOpts {
	return &queryTransactionOpts{
		includeProof: false,
		merkleFormat: MerkleFormatTSC,
	}
}

// WithQueryMerkleProof will request that a merkle proof is returned in the query response.
func WithQueryMerkleProof() QueryTransactionOptFunc {
	return func(o *queryTransactionOpts) {
		o.includeProof = true
	}
}

// WithoutQueryMerkleProof is the default option and doesn't need to be passed however can be
// added for code clarity.
func WithoutQueryMerkleProof() QueryTransactionOptFunc {
	return func(o *queryTransactionOpts) {
		o.includeProof = false
	}
}

// WithQueryMerkleFormat can be used to overwrite the default format of TSC with another value.
func WithQueryMerkleFormat(s string) QueryTransactionOptFunc {
	return func(o *queryTransactionOpts) {
		o.merkleFormat = s
	}
}

// WithQueryTSCMerkleFormat is the default option and doesn't need to be passed however can be
// added for code clarity.
func WithQueryTSCMerkleFormat() QueryTransactionOptFunc {
	return func(o *queryTransactionOpts) {
		o.merkleFormat = MerkleFormatTSC
	}
}

// QueryTransaction will fire a Merchant API request to check the status of a transaction
//
// This endpoint is used to check the current status of a previously submitted transaction.
// It returns a JSONEnvelope with a payload that contains the transaction status.
// The purpose of the envelope is to ensure strict consistency in the message content for
// the purpose of signing responses.
//
// You can provide optional arguments using the WithQuery... option functions, an example is shown:
//
//	QueryTransaction(ctx, miner, abc123, WithQueryMerkleProof(), WithQueryTSCMerkleFormat())
//
// This is backwards compatible with the previous non-optional version of this function and be called with 0 options:
//
//	QueryTransaction(ctx, miner, abc123)
//
// In this case the defaults are used which is to not request a proof.
// Specs: https://github.com/bitcoin-sv-specs/brfc-merchantapi#4-query-transaction-status
func (c *Client) QueryTransaction(ctx context.Context, miner *Miner, txID string, opts ...QueryTransactionOptFunc) (*QueryTransactionResponse, error) {

	// Make sure we have a valid miner
	if miner == nil {
		return nil, errors.New("miner was nil")
	}

	// Make the HTTP request
	result := queryTransaction(ctx, c, miner, txID, opts...)
	if result.Response.Error != nil {
		return nil, result.Response.Error
	}

	queryResponse := &QueryTransactionResponse{
		JSONEnvelope: JSONEnvelope{
			ApiType: c.apiType,
			Miner:   result.Miner,
		},
	}

	var modelAdapter QueryTxModelAdapter

	switch c.apiType {
	case MAPI:
		model := &mapi.QueryTxModel{}
		err := queryResponse.process(result.Miner, result.Response.BodyContents)
		if err != nil || len(queryResponse.Payload) <= 0 {
			return nil, err
		}

		err = json.Unmarshal([]byte(queryResponse.Payload), model)
		if err != nil {
			return nil, err
		}

		modelAdapter = &QueryTxMapiAdapter{QueryTxModel: model}

	case Arc:
		model := &arc.QueryTxModel{}
		err := json.Unmarshal(result.Response.BodyContents, model)
		if err != nil {
			return nil, err
		}

		modelAdapter = &QueryTxArcAdapter{QueryTxModel: model}

	default:
		return nil, errors.New("unsupported api type")
	}

	queryResponse.Query = modelAdapter.GetQueryTxResponse()

	// Valid?
	if queryResponse.Query == nil {
		return nil, errors.New("failed getting query response from: " + miner.Name)
	}

	isValid, err := queryResponse.IsValid()
	if err != nil {
		return nil, err
	}

	queryResponse.Validated = isValid

	// Return the fully parsed response
	return queryResponse, nil
}

// queryTransaction will fire the HTTP request to retrieve the tx status
func queryTransaction(ctx context.Context, client *Client, miner *Miner, txHash string, opts ...QueryTransactionOptFunc) (result *internalResult) {
	defaultOpts := defaultQueryOpts()

	// overwrite defaults with any provided options
	for _, o := range opts {
		o(defaultOpts)
	}
	sb := strings.Builder{}

	api, err := MinerAPIByMinerID(client.minerAPIs, miner.MinerID, client.apiType)
	if err != nil {
		result.Response = &RequestResponse{Error: err}
		return
	}

	route, err := ActionRouteByAPIType(QueryTx, client.apiType)
	if err != nil {
		result.Response = &RequestResponse{Error: err}
		return
	}

	sb.WriteString(api.URL + route + txHash)
	if defaultOpts.includeProof {
		sb.WriteString("?merkleProof=true&merkleFormat=" + defaultOpts.merkleFormat)
	}
	result = &internalResult{Miner: miner}
	queryURL, err := url.Parse(sb.String())
	if err != nil {
		result.Response = &RequestResponse{Error: err}
		return
	}

	result.Response = httpRequest(ctx, client, &httpPayload{
		Method: http.MethodGet,
		URL:    queryURL.String(),
		Token:  api.Token,
	})
	return
}

func (m *QueryTxMapiAdapter) GetQueryTxResponse() *QueryTxResponse {
	response := &QueryTxResponse{
		TxID:        m.TxID,
		BlockHash:   m.BlockHash,
		Timestamp:   m.Timestamp,
		BlockHeight: m.BlockHeight,
	}

	// Fields specific to mAPI
	response.APIVersion = m.APIVersion
	response.ReturnResult = m.ReturnResult
	response.ResultDescription = m.ResultDescription
	response.MinerID = m.MinerID
	response.Confirmations = m.Confirmations
	response.TxSecondMempoolExpiry = m.TxSecondMempoolExpiry
	response.MerkleProof = m.MerkleProof

	return response
}

func (m *QueryTxArcAdapter) GetQueryTxResponse() *QueryTxResponse {
	response := &QueryTxResponse{
		TxID:        m.TxID,
		BlockHash:   m.BlockHash,
		Timestamp:   m.Timestamp.String(),
		BlockHeight: m.BlockHeight,
	}

	// Fields specific to ArcAPI
	response.TxStatus = m.TxStatus

	return response
}
