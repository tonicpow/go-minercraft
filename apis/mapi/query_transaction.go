package mapi

import "github.com/libsv/go-bc"

type QueryTxModel struct {
	APIVersion            string          `json:"apiVersion"`
	Timestamp             string          `json:"timestamp"`
	TxID                  string          `json:"txid"`
	ReturnResult          string          `json:"returnResult"`
	ResultDescription     string          `json:"resultDescription"`
	BlockHash             string          `json:"blockHash"`
	BlockHeight           int64           `json:"blockHeight"`
	MinerID               string          `json:"minerId"`
	Confirmations         int64           `json:"confirmations"`
	TxSecondMempoolExpiry int64           `json:"txSecondMempoolExpiry"`
	MerkleProof           *bc.MerkleProof `json:"merkleProof"`
}
